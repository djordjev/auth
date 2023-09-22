package main

import (
	"fmt"
	"net/http"

	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/utils"
	"github.com/djordjev/auth/packages/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Read config
	config, err := utils.BuildConfigFromEnv()
	if err != nil {
		panic(err)
	}

	// Auto migrate database
	dbUrl := config.GetConnectionString()
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := models.AutoMigrate(db); err != nil {
		panic(err)
	}

	// Start up
	r := http.NewServeMux()
	server := server.NewServer(r, config)
	server.Mount("/")

	fmt.Printf("Running server on port %s\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}
