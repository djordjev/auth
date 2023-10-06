package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/djordjev/pg-mig/migrations"
)

func AutoMigrate(connectionString string) error {
	envMigrations := os.Getenv("TEST_DB_MIGRATIONS")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var migrationsFolder string
	if envMigrations != "" {
		migrationsFolder = envMigrations
	} else {
		migrationsFolder = filepath.Join(wd, "../../migrations")
	}

	fmt.Println("👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿👿", migrationsFolder)

	trimProtocol := strings.TrimPrefix(connectionString, "postgres://")

	credentialsArr := strings.Split(trimProtocol, "@")
	credentials := credentialsArr[0]

	dbInfo := strings.Split(credentialsArr[1], "/")
	dbName := dbInfo[1]

	hostPort := strings.Split(dbInfo[0], ":")

	host := hostPort[0]

	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return err
	}

	runner := migrations.NewRunner(
		host,
		credentials,
		dbName,
		port,
		migrationsFolder,
	)

	return runner.Run([]string{})
}
