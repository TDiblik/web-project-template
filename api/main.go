package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/TDiblik/project-template/api/constants"
	"github.com/TDiblik/project-template/api/database"
	"github.com/TDiblik/project-template/api/router"
	"github.com/TDiblik/project-template/api/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	utils.SetupENV()
	utils.SetupValidator()

	if !fiber.IsChild() {
		utils.Logger.Info("Initializing the Template API server")
	}

	if !fiber.IsChild() {
		utils.Logger.Info("Checking database connectivity: start")
	}
	db, err := database.CreateConnection()
	if err != nil {
		utils.Logger.Fatalw("Unable to connect to a database", "error", err)
	}
	if err := db.Ping(); err != nil {
		utils.Logger.Fatalw("Unable to ping a database", "error", err)
	}
	if !fiber.IsChild() {
		utils.Logger.Info("Checking database connectivity: done")
	}

	if !fiber.IsChild() {
		utils.Logger.Info("Running database migrations: start")
		m, err := migrate.New(utils.EnvData.DB_MIGRATIONS_PATH, utils.EnvData.DB_CONNECTION_STRING)
		if err != nil {
			utils.Logger.Fatalw("Unable to run database migrations (when creating migrate instance)", "error", err)
		}
		if utils.EnvData.DB_DEV_FORCE_MIGRATE_DOWN {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				utils.Logger.Fatalw("Unable to run database migrations (down)", "error", err)
			} else if err == migrate.ErrNoChange {
				utils.Logger.Info("No new migrations to run (down)")
			} else {
				utils.Logger.Info("Successfully run migrations (down)")
			}
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			utils.Logger.Fatalw("Unable to run database migrations (up)", "error", err)
		} else if err == migrate.ErrNoChange {
			utils.Logger.Info("No new migrations to run (up)")
		} else {
			utils.Logger.Info("Successfully run migrations (up)")
		}
		utils.Logger.Info("Running database migrations: done")
	}

	if !fiber.IsChild() {
		utils.Logger.Info("Setting up the app: start")
	}
	app := fiber.New(fiber.Config{
		AppName:       "Template APP",
		CaseSensitive: true,
		BodyLimit:     15 * 1024 * 1024, // 15mb
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   120 * time.Second,
	})

	app.Use(recover.New(
		recover.Config{
			Next:              recover.ConfigDefault.Next,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
			EnableStackTrace:  utils.EnvData.Debug,
		},
	))

	loggerConfig := logger.ConfigDefault
	loggerConfig.Next = func(c fiber.Ctx) bool {
		path := c.Path()
		return path == "/api/health" || strings.HasPrefix(path, "/assets/")
	}
	if !utils.EnvData.Debug {
		loggerConfig.Format = `{"time":"${time}","status":${status},"latency":"${latency}","ip":"${ip}","method":"${method}","path":"${path}","error":"${error}"}` + "\n"
	}
	app.Use(logger.New(loggerConfig))

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{utils.EnvData.FE_PROD_URL}, // frontend origin
			AllowCredentials: true,
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", constants.TOKEN_HEADER_NAME},
		},
	))

	app.Use(helmet.New())

	if !fiber.IsChild() {
		utils.Logger.Info("Setting up the app: done")
		utils.Logger.Info("Setting up the routes: start")
	}
	router.SetupRoutes(app)
	if !fiber.IsChild() {
		utils.Logger.Info("Setting up the routes: done")
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if !fiber.IsChild() {
			utils.Logger.Info("We received an interrupt signal, shutting down server...")
		}
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			utils.Logger.Errorw("Server shutdown failed", "error", err)
		}
		if err := db.Close(); err != nil {
			utils.Logger.Errorw("Database connection close failed", "error", err)
		}
		close(idleConnsClosed)
	}()

	cron_ctx := utils.WithSignalCancel("cron jobs")
	utils.SetupCronJobs(cron_ctx)

	if !fiber.IsChild() {
		utils.Logger.Info("Initialization completed")
	}

	if err := app.Listen(":"+utils.EnvData.API_PORT, fiber.ListenConfig{
		EnablePrefork:     !utils.EnvData.Debug,
		EnablePrintRoutes: false,
	}); err != nil {
		utils.Logger.Fatalw("Error starting server", "error", err)
	}
	<-idleConnsClosed

	utils.Logger.Info("Server successfully shut down")
}
