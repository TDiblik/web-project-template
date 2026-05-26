package utils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
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

func SetupENV(env_files ...string) {
	LogIfMaster("Setting up env variables: start")

	err := godotenv.Load(env_files...)
	if err != nil {
		LogIfMaster("Unable to load .env file: ", err)
		LogIfMaster("This is normal in production environments, since all environment variables are set in the cloud.")
	}

	childStatus := "master"
	if fiber.IsChild() {
		childStatus = "child"
	}
	switch getEnvKeyOrPanic("GO_ENV") {
	case "development":
		EnvData.Debug = true
		log.SetPrefix(fmt.Sprintf("[DEBUG] - %d (%s) - ", os.Getpid(), childStatus) + log.Prefix())
	case "production":
		EnvData.Debug = false
		log.SetPrefix(fmt.Sprintf("[PROD] - %d (%s) - ", os.Getpid(), childStatus) + log.Prefix())
	default:
		log.Fatalln("Error determening GO_ENV (", os.Getenv("GO_ENV"), ")")
	}

	EnvData.API_PORT = getEnvKeyOrPanic("API_PORT")
	EnvData.DB_CONNECTION_STRING = getEnvKeyOrPanic("DB_CONNECTION_STRING")
	EnvData.DB_MIGRATIONS_PATH = getEnvKeyOrPanic("DB_MIGRATIONS_PATH")
	if strings.ToLower(getEnvKeyOrPanic("DB_DEV_FORCE_MIGRATE_DOWN")) == "true" {
		if !EnvData.Debug {
			log.Fatalln("Cannot use DB_DEV_FORCE_MIGRATE_DOWN while in production mode!")
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
		log.Fatal("Error ensuring images path (\"", EnvData.IMAGES_PATH, "\"): ", err)
	}
	EnvData.IMAGES_PATH_AVATAR = GetAvatarImageFolder()
	if err := os.MkdirAll(EnvData.IMAGES_PATH_AVATAR, FoldrePerms); err != nil {
		log.Fatal("Error ensuring images path (\"", EnvData.IMAGES_PATH_AVATAR, "\"): ", err)
	}
	EnvData.IMAGES_PATH_TEMP = GetTempImageFolder()
	if err := os.MkdirAll(EnvData.IMAGES_PATH_TEMP, FoldrePerms); err != nil {
		log.Fatal("Error ensuring images path (\"", EnvData.IMAGES_PATH_TEMP, "\"): ", err)
	}

	log.Println("Setting up env variables: done")
}

func getEnvKeyOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		log.Fatal("Error loading ", key)
	}
	return val
}

func Log(v ...any) {
	log.Println(v...)
}

func LogErr(e error) {
	log.Println("ERROR: ", e)
}

func LogIfMaster(v ...any) {
	if !fiber.IsChild() {
		Log(v...)
	}
}
