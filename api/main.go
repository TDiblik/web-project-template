package main

import (
	"log"

	"github.com/TDiblik/project-template/api/constants"
	"github.com/TDiblik/project-template/api/database"
	"github.com/TDiblik/project-template/api/router"
	"github.com/TDiblik/project-template/api/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	utils.LogIfMaster("Initializing the Template API server")

	utils.SetupENV()
	utils.SetupValidator()

	log.Println("Checking database connectivity: start")
	db, err := database.CreateConnection()
	if err != nil {
		log.Fatalln("Unable to connect to a database", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln("Unable to ping a database", err)
	}
	log.Println("Checking database connectivity: done")

	if !fiber.IsChild() {
		log.Println("Running database migrations: start")
		m, err := migrate.New(utils.EnvData.DB_MIGRATIONS_PATH, utils.EnvData.DB_CONNECTION_STRING)
		if err != nil {
			log.Fatalln("Unable to run database migrations (when creating migrate instance)", err)
		}
		if utils.EnvData.DB_DEV_FORCE_MIGRATE_DOWN {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalln("Unable to run database migrations (down)", err)
			} else if err == migrate.ErrNoChange {
				log.Println("No new migrations to run (down)")
			} else {
				log.Println("Successfully run migrations (down)")
			}
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalln("Unable to run database migrations (up)", err)
		} else if err == migrate.ErrNoChange {
			log.Println("No new migrations to run (up)")
		} else {
			log.Println("Successfully run migrations (up)")
		}
		log.Println("Running database migrations: done")
	}

	log.Println("Setting up the app: start")
	app := fiber.New(fiber.Config{
		AppName:       "Template APP",
		CaseSensitive: true,
		BodyLimit:     15 * 1024 * 1024, // 15mb
	})
	app.Use(logger.New())
	app.Use(recover.New(
		recover.Config{
			Next:              recover.ConfigDefault.Next,
			StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
			EnableStackTrace:  utils.EnvData.Debug,
		},
	))
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{utils.EnvData.FE_PROD_URL}, // frontend origin
			AllowCredentials: true,
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", constants.TOKEN_HEADER_NAME},
		},
	))
	log.Println("Setting up the app: done")

	log.Println("Setting up the routes: start")
	router.SetupRoutes(app)
	log.Println("Setting up the routes: done")

	cron_ctx := utils.WithSignalCancel("cron jobs")
	utils.SetupCronJobs(cron_ctx)
	log.Println("Initialization completed")
	log.Fatal(app.Listen(":"+utils.EnvData.API_PORT, fiber.ListenConfig{
		EnablePrefork:     !utils.EnvData.Debug,
		EnablePrintRoutes: false,
	}))
}
