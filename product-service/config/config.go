package config

import (
	"os"
	"strconv"
)

// Config contains the configuration for the service
type Config struct {
	Port                int
	UsersServiceBaseUrl string
	DBUsername          string
	DBPassword          string
	DBHost              string
	DBPort              int
	Database            string
	SecretKey           string
}

// LoadConfig returns the config from the environment variables
func LoadConfig() Config {

	config := Config{}

	if _, ok := os.LookupEnv("PORT"); ok {
		portString := os.Getenv("PORT")
		port, err := strconv.ParseInt(portString, 10, 0)
		if err == nil {
			config.Port = int(port)
		}
	}

	if _, ok := os.LookupEnv("PRODUCTS_SERVICE_URL"); ok {
		config.UsersServiceBaseUrl = os.Getenv("PRODUCTS_SERVICE_URL")
	}

	if _, ok := os.LookupEnv("DB_USERNAME"); ok {
		config.DBUsername = os.Getenv("DB_USERNAME")
	}

	if _, ok := os.LookupEnv("DB_PASSWORD"); ok {
		config.DBPassword = os.Getenv("DB_PASSWORD")
	}

	if _, ok := os.LookupEnv("DB_HOST"); ok {
		config.DBHost = os.Getenv("DB_HOST")
	}

	if _, ok := os.LookupEnv("DB_PORT"); ok {
		portString := os.Getenv("DB_PORT")
		port, err := strconv.Atoi(portString)
		if err == nil {
			config.DBPort = port
		}
	}

	if _, ok := os.LookupEnv("DB"); ok {
		config.Database = os.Getenv("DB")
	}

	if _, ok := os.LookupEnv("SECRET_KEY"); ok {
		config.SecretKey = os.Getenv("SECRET_KEY")
	}

	return config
}
