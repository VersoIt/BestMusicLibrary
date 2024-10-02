package cfg

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	DbHost               string
	DbPort               string
	DbUser               string
	DbPassword           string
	DbName               string
	DbSSLMode            string
	ExternalApiClientUrl string
	ServerPort           string
}

var (
	config Config
	once   sync.Once
)

func Get() Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			logrus.Error("error loading .env file")
		}
		config = Config{}
		config.DbHost = os.Getenv("DB_HOST")
		config.DbPort = os.Getenv("DB_PORT")
		config.DbUser = os.Getenv("DB_USER")
		config.DbPassword = os.Getenv("DB_PASSWORD")
		config.DbName = os.Getenv("DB_NAME")
		config.ExternalApiClientUrl = os.Getenv("EXTERNAL_API_CLIENT_URL")
		config.ServerPort = os.Getenv("SERVER_PORT")
		config.DbSSLMode = os.Getenv("DB_SSL_MODE")
	})

	return config
}
