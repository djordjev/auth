package models

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	FailedToConnect = iota + 1000
	FailedToPing
	FailedToCloseConnection
	FailedToRunMigrations
)

var dbConnection *pgxpool.Pool

func getTestConnectionString() string {
	envDBUrl := os.Getenv("TEST_DB_URL")

	if envDBUrl != "" {
		return envDBUrl
	} else {
		return "postgres://tester:testee@localhost:5433/testdb"
	}
}

func TestMain(m *testing.M) {
	var err error

	count := 10
	dbUrl := getTestConnectionString()
	fmt.Println("💿 Connecting to the database")

	for i := 0; i < count; i++ {
		dbConnection, err = pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			if i == count-1 {
				os.Exit(FailedToConnect)
			}

			// failed to connect to db
			waitTime := time.Duration(5) * time.Second
			fmt.Printf("♻️  Failed to connect to database %s, attempting again after %s seconds\n", dbUrl, waitTime)
			time.Sleep(waitTime)
		} else {
			break
		}
	}

	if err = dbConnection.Ping(context.Background()); err != nil {
		fmt.Println("🛑 Failed to ping database")
		os.Exit(FailedToPing)
	}

	fmt.Println("🎉🎉🎉 connected to database")

	defer dbConnection.Close()

	if err = AutoMigrate(dbUrl); err != nil {
		fmt.Println("🛑 Failed to run migrations")
		os.Exit(FailedToRunMigrations)
	}

	res := m.Run()

	os.Exit(res)
}
