package utils

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/spotify"
)

type IEnvData struct {
	Debug bool

	API_PORT                  string
	API_PROD_URL              string
	FE_PROD_URL               string
	DB_CONNECTION_STRING      string
	DB_MIGRATIONS_PATH        string
	DB_DEV_FORCE_MIGRATE_DOWN bool

	AUTH_JWT_SECRET   string
	AUTH_SECRET_BYTES []byte

	OAUTH_JWT_SECRET   string
	OAUTH_SECRET_BYTES []byte

	OAUTH_GITHUB_CLIENT_ID     string
	OAUTH_GITHUB_CLIENT_SECRET string
	OAUTH_GITHUB_CONFIG        *oauth2.Config

	OAUTH_GOOGLE_CLIENT_ID     string
	OAUTH_GOOGLE_CLIENT_SECRET string
	OAUTH_GOOGLE_CONFIG        *oauth2.Config

	OAUTH_FACEBOOK_CLIENT_ID     string
	OAUTH_FACEBOOK_CLIENT_SECRET string
	OAUTH_FACEBOOK_CONFIG        *oauth2.Config

	OAUTH_SPOTIFY_CLIENT_ID     string
	OAUTH_SPOTIFY_CLIENT_SECRET string
	OAUTH_SPOTIFY_CONFIG        *oauth2.Config

	IMAGES_PATH        string
	IMAGES_PATH_AVATAR string
	IMAGES_PATH_TEMP   string
}

var EnvData IEnvData
var FoldrePerms fs.FileMode = 0o777

var Logger *zap.SugaredLogger

func initLogger(debug bool, childStatus string) {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}
	config.InitialFields = map[string]interface{}{
		"pid":  os.Getpid(),
		"proc": childStatus,
	}

	zapLogger, err := config.Build()
	if err != nil {
		fmt.Printf("Failed to initialize zap logger: %v\n", err)
		os.Exit(1)
	}
	Logger = zapLogger.Sugar()
}

func SetupENV(env_files ...string) {
	err := godotenv.Load(env_files...)

	childStatus := "master"
	if fiber.IsChild() {
		childStatus = "child"
	}

	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		fmt.Println("Error determining GO_ENV (missing)")
		os.Exit(1)
	}

	initLogger(goEnv == "development", childStatus)

	if err != nil {
		if !fiber.IsChild() && Logger != nil {
			Logger.Errorw("Unable to load .env file", "error", err)
			Logger.Info("This is normal in production environments, since all environment variables are set in the cloud.")
		}
	}

	if !fiber.IsChild() && Logger != nil {
		Logger.Info("Setting up env variables: start")
	}

	switch goEnv {
	case "development":
		EnvData.Debug = true
	case "production":
		EnvData.Debug = false
	default:
		Logger.Fatalf("Error determening GO_ENV (%s)", goEnv)
	}

	EnvData.API_PORT = getEnvKeyOrPanic("API_PORT")
	EnvData.DB_CONNECTION_STRING = getEnvKeyOrPanic("DB_CONNECTION_STRING")
	EnvData.DB_MIGRATIONS_PATH = getEnvKeyOrPanic("DB_MIGRATIONS_PATH")
	if strings.ToLower(getEnvKeyOrPanic("DB_DEV_FORCE_MIGRATE_DOWN")) == "true" {
		if !EnvData.Debug {
			Logger.Fatal("Cannot use DB_DEV_FORCE_MIGRATE_DOWN while in production mode!")
		}
		EnvData.DB_DEV_FORCE_MIGRATE_DOWN = true
	} else {
		EnvData.DB_DEV_FORCE_MIGRATE_DOWN = false
	}
	EnvData.API_PROD_URL = getEnvKeyOrPanic("API_PROD_URL")
	if !strings.HasSuffix(EnvData.API_PROD_URL, "/") {
		EnvData.API_PROD_URL += "/"
	}
	EnvData.FE_PROD_URL = getEnvKeyOrPanic("FE_PROD_URL")
	if !strings.HasSuffix(EnvData.FE_PROD_URL, "/") {
		EnvData.FE_PROD_URL += "/"
	}

	EnvData.AUTH_JWT_SECRET = getEnvKeyOrPanic("AUTH_JWT_SECRET")
	EnvData.AUTH_SECRET_BYTES = []byte(EnvData.AUTH_JWT_SECRET)

	EnvData.OAUTH_JWT_SECRET = getEnvKeyOrPanic("OAUTH_JWT_SECRET")
	EnvData.OAUTH_SECRET_BYTES = []byte(EnvData.OAUTH_JWT_SECRET)

	// when adding a new oauth provider and user table fields, add the checks here:
	EnvData.OAUTH_GITHUB_CLIENT_ID = getEnvKeyOrPanic("OAUTH_GITHUB_CLIENT_ID")
	EnvData.OAUTH_GITHUB_CLIENT_SECRET = getEnvKeyOrPanic("OAUTH_GITHUB_CLIENT_SECRET")
	EnvData.OAUTH_GITHUB_CONFIG = &oauth2.Config{
		ClientID:     EnvData.OAUTH_GITHUB_CLIENT_ID,
		ClientSecret: EnvData.OAUTH_GITHUB_CLIENT_SECRET,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  JoinUrlOrPanic(EnvData.FE_PROD_URL, "/login/oauth/redirect"),
	}

	EnvData.OAUTH_GOOGLE_CLIENT_ID = getEnvKeyOrPanic("OAUTH_GOOGLE_CLIENT_ID")
	EnvData.OAUTH_GOOGLE_CLIENT_SECRET = getEnvKeyOrPanic("OAUTH_GOOGLE_CLIENT_SECRET")
	EnvData.OAUTH_GOOGLE_CONFIG = &oauth2.Config{
		ClientID:     EnvData.OAUTH_GOOGLE_CLIENT_ID,
		ClientSecret: EnvData.OAUTH_GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  JoinUrlOrPanic(EnvData.FE_PROD_URL, "/login/oauth/redirect"),
	}

	EnvData.OAUTH_FACEBOOK_CLIENT_ID = getEnvKeyOrPanic("OAUTH_FACEBOOK_CLIENT_ID")
	EnvData.OAUTH_FACEBOOK_CLIENT_SECRET = getEnvKeyOrPanic("OAUTH_FACEBOOK_CLIENT_SECRET")
	EnvData.OAUTH_FACEBOOK_CONFIG = &oauth2.Config{
		ClientID:     EnvData.OAUTH_FACEBOOK_CLIENT_ID,
		ClientSecret: EnvData.OAUTH_FACEBOOK_CLIENT_SECRET,
		Scopes:       []string{"public_profile", "email", "user_link"},
		Endpoint:     facebook.Endpoint,
		RedirectURL:  JoinUrlOrPanic(EnvData.FE_PROD_URL, "/login/oauth/redirect"),
	}

	EnvData.OAUTH_SPOTIFY_CLIENT_ID = getEnvKeyOrPanic("OAUTH_SPOTIFY_CLIENT_ID")
	EnvData.OAUTH_SPOTIFY_CLIENT_SECRET = getEnvKeyOrPanic("OAUTH_SPOTIFY_CLIENT_SECRET")
	EnvData.OAUTH_SPOTIFY_CONFIG = &oauth2.Config{
		ClientID:     EnvData.OAUTH_SPOTIFY_CLIENT_ID,
		ClientSecret: EnvData.OAUTH_SPOTIFY_CLIENT_SECRET,
		Scopes:       []string{"user-read-email", "user-read-private"},
		Endpoint:     spotify.Endpoint,
		RedirectURL:  JoinUrlOrPanic(EnvData.FE_PROD_URL, "/login/oauth/redirect"),
	}

	EnvData.IMAGES_PATH = getEnvKeyOrPanic("IMAGES_PATH")
	if err := os.MkdirAll(EnvData.IMAGES_PATH, FoldrePerms); err != nil {
		Logger.Fatalw("Error ensuring images path", "path", EnvData.IMAGES_PATH, "error", err)
	}
	EnvData.IMAGES_PATH_AVATAR = GetAvatarImageFolder()
	if err := os.MkdirAll(EnvData.IMAGES_PATH_AVATAR, FoldrePerms); err != nil {
		Logger.Fatalw("Error ensuring images path", "path", EnvData.IMAGES_PATH_AVATAR, "error", err)
	}
	EnvData.IMAGES_PATH_TEMP = GetTempImageFolder()
	if err := os.MkdirAll(EnvData.IMAGES_PATH_TEMP, FoldrePerms); err != nil {
		Logger.Fatalw("Error ensuring images path", "path", EnvData.IMAGES_PATH_TEMP, "error", err)
	}

	Logger.Info("Setting up env variables: done")
}

func getEnvKeyOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		Logger.Fatalf("Error loading %s", key)
	}
	return val
}
