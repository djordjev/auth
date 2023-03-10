package utils

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost string
	DBName string
	DBPass string
	DBPort uint
	DBUser string
	Port   string
}

func BuildConfigFromEnv() (Config, error) {
	config := Config{}

	config.DBHost = os.Getenv("DB_HOST")
	config.DBPass = os.Getenv("DB_PASS")
	config.DBUser = os.Getenv("DB_USER")
	config.DBName = os.Getenv("DB_NAME")
	config.Port = os.Getenv("PORT")

	if os.Getenv("DB_PORT") == "" {
		// If not set default to PostgreSQL default 5432 port
		config.DBPort = 5432
	} else {
		// Otherwise parse DB_PORT env
		dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			return Config{}, err
		}
		config.DBPort = uint(dbPort)
	}

	return config, nil
}

func (config Config) GetConnectionString() string {
	if config.DBName == "" || config.DBUser == "" || config.DBPass == "" || config.DBHost == "" {
		panic(fmt.Errorf("missing database fields in configuration"))
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
}
