package utils

import (
	"fmt"
	"os"
	"strconv"
)

type Mailjet struct {
	ApiKey    string
	SecretKey string
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
	Mailjet             Mailjet
	VerificationLink    string
	ForgetPasswordLink  string
	Sender              string
	RedisPort           uint
	RedisHost           string
	RedisPassword       string
	RedisDatabase       int
	SessionCookie       string
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

	config.Mailjet.ApiKey = os.Getenv("MAILJET_API_KEY")
	config.Mailjet.SecretKey = os.Getenv("MAILJET_SECRET_KEY")

	config.VerificationLink = os.Getenv("VERIFICATION_LINK")
	config.ForgetPasswordLink = os.Getenv("FORGET_PASSWORD_LINK")
	config.Sender = os.Getenv("SENDER")

	if os.Getenv("REDIS_PORT") == "" {
		config.RedisPort = 6379
	} else {
		redisPort, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 64)
		if err != nil {
			return Config{}, err
		}

		config.RedisPort = uint(redisPort)
	}

	config.RedisHost = os.Getenv("REDIS_HOST")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")

	if os.Getenv("REDIS_DB") == "" {
		config.RedisDatabase = 0
	} else {
		redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			return Config{}, err
		}

		config.RedisDatabase = redisDB
	}

	config.SessionCookie = os.Getenv("SESSION_COOKIE")
	if config.SessionCookie == "" {
		config.SessionCookie = "_tkn"
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
	return config.Mailjet.ApiKey != "" && config.Mailjet.SecretKey != ""
}
