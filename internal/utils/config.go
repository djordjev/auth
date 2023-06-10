package utils

import (
	"fmt"
	"os"
	"strconv"
)

type Smtp struct {
	Url      string
	Password string
	Port     int
	Username string
}

type Config struct {
	DBHost              string
	DBName              string
	DBPass              string
	DBPort              uint
	DBUser              string
	Domain              string
	Port                string
	GoEnv               string
	RequireVerification bool
	Smtp                Smtp
}

func BuildConfigFromEnv() (Config, error) {
	config := Config{}

	config.DBHost = os.Getenv("DB_HOST")
	config.DBPass = os.Getenv("DB_PASS")
	config.DBUser = os.Getenv("DB_USER")
	config.DBName = os.Getenv("DB_NAME")
	config.Port = os.Getenv("PORT")
	config.GoEnv = os.Getenv("GO_ENV")
	config.Domain = os.Getenv("DOMAIN")

	if config.GoEnv == "" {
		config.GoEnv = "development"
	}

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

	if os.Getenv("REQUIRE_VERIFICATION") == "true" {
		config.RequireVerification = true
	} else {
		config.RequireVerification = false
	}

	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpURI := os.Getenv("SMTP_URI")

	if smtpUsername != "" && smtpPassword != "" && smtpURI != "" {
		smtp := Smtp{
			Url:      smtpURI,
			Password: smtpPassword,
			Username: smtpUsername,
		}

		if smtpPort == "" {
			smtpPort = "587"
		}

		if smtpPortNum, err := strconv.Atoi(smtpPort); err == nil {
			smtp.Port = smtpPortNum
		} else {
			return config, err
		}

		config.Smtp = smtp
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

func (config Config) IsDev() bool {
	return config.GoEnv == "development"
}

func (config Config) HasEmailSetup() bool {
	return config.Smtp.Url != "" && config.Smtp.Username != "" && config.Smtp.Password != "" && config.Smtp.Port != 0
}
