package config

import (
	"account-test/postgres"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB *postgres.DBConfig
}

func InitReader() {
	environment := ""
	if len(os.Args) < 2 {
		log.Fatalf("Env not supplied in argument")
	} else {
		environment = os.Args[1]
	}

	err := godotenv.Load(environment + ".env")
	if err != nil {
		log.Fatalf("Error loading %s.env file", environment)
	}
}

func Init() AppConfig {

	appConfig := AppConfig{
		DB: &postgres.DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Schema:   os.Getenv("DB_SCHEMA"),
		},
	}

	return appConfig
}
