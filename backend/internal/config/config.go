package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	MongoConfig struct {
		ConnectionURI string `env:"URI" envDefault:"mongodb://localhost:27017"`
		Database      string `env:"DB" envDefault:"bar"`
		Timeout       uint64 `env:"TIMEOUT" envDefault:"30"`
	} `envPrefix:"BAR_MONGO_"`

	ApiConfig struct {
		SessionSecret      string `env:"SESSION_SECRET"`
		AdminSessionSecret string `env:"ADMIN_SESSION_SECRET"`
		Port               string `env:"PORT" envDefault:":8080"`
		BasePath           string `env:"BASE_PATH" envDefault:"http://localhost:8080"`
		LocalToken         string `env:"LOCAL_TOKEN"`
	} `envPrefix:"BAR_API_"`

	OauthConfig struct {
		GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	} `envPrefix:"BAR_OAUTH_"`

	StorageConfig struct {
		StoragePath string `env:"STORAGE_PATH" envDefault:"./storage"`
	} `envPrefix:"BAR_STORAGE_"`

	LogLevel string `env:"BAR_LOG_LEVEL" envDefault:"info"`
}

var config Config

func GetConfig() Config {
	return config
}

func init() {
	godotenv.Load()
	if err := env.Parse(&config); err != nil {
		logrus.Fatal(err)
	}

	logrus.SetLevel(logrus.InfoLevel)
	if config.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Loaded config: ", fmt.Sprintf("%+v", config))
}
