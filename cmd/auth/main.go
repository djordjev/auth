package main

import (
	"fmt"
	"net/http"

	"github.com/djordjev/auth/internal/api"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/notify"
	"github.com/djordjev/auth/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := utils.BuildConfigFromEnv()
	if err != nil {
		panic(err)
	}

	r := http.NewServeMux()

	// Init database
	dbUrl := config.GetConnectionString()
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	if err = models.AutoMigrate(db); err != nil {
		panic(err)
	}

	// Init redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDatabase,
	})

	repo := models.NewRepository(db, client)

	// Init api
	logger := utils.MustBuildLogger(config)
	var notifier domain.Notifier
	if config.IsDev() {
		notifier = notify.SilentNotifier{}
	} else {
		notifier = notify.NewMailjetNotifier(config)
	}

	// Init app domain
	appDomain := domain.NewDomain(repo, config, notifier)

	appApi := api.NewApi(config, r, appDomain, logger)

	// Start up
	appApi.Mount("/")

	fmt.Println("running")
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}
