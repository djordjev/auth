package main

import (
	"fmt"
	"net/http"

	"github.com/djordjev/auth/internal/utils"
	"github.com/djordjev/auth/packages/server"
)

func main() {
	// Read config
	config, err := utils.BuildConfigFromEnv()
	if err != nil {
		panic(err)
	}

	// Auto migrate database
	// dbUrl := config.GetConnectionString()
	// pool, err := pgxpool.New(context.Background(), dbUrl)
	// if err != nil {
	// 	panic(err)
	// }

	// defer pool.Close()

	// TODO: auto migrate
	// if err := models.AutoMigrate(db); err != nil {
	// 	panic(err)
	// }

	// Start up
	r := http.NewServeMux()
	server := server.NewServer(r, config)

	defer server.Close()

	server.Mount("/")

	fmt.Printf("Running server on port %s\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}
