package models

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	FailedToConnect = iota + 1000
	FailedToPing
	FailedToCloseConnection
	FailedToRunMigrations
)

var dbConnection *gorm.DB

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

	count := 7
	dbUrl := getTestConnectionString()
	fmt.Println("💿 Connecting to the database")

	for i := 0; i < count; i++ {
		dbConnection, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			if i == count-1 {
				os.Exit(FailedToConnect)
			}

			// failed to connect to db
			waitTime := time.Duration(i+1) * 5 * time.Second
			fmt.Printf("♻️  Failed to connect to database %s, attempting again after %s seconds\n", dbUrl, waitTime)
			time.Sleep(waitTime)
		} else {
			break
		}
	}

	cnn, _ := dbConnection.DB()
	err = cnn.Ping()
	if err != nil {
		fmt.Println("🛑 Failed to ping database")
		os.Exit(FailedToPing)
	}

	fmt.Println("🎉🎉🎉 connected to database")

	if err = AutoMigrate(dbConnection); err != nil {
		fmt.Println("🛑 Failed to run migrations")
		os.Exit(FailedToRunMigrations)
	}

	defer func() {
		err = cnn.Close()

		if err != nil {
			fmt.Println("🛑 Failed to close connection")
			os.Exit(FailedToCloseConnection)
		}
	}()

	res := m.Run()

	os.Exit(res)
}
