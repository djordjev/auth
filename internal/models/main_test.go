package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

const (
	FailedToConnect = iota + 1000
	FailedToPing
	FailedToCloseConnection
	FailedToRunMigrations
)

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
	var db *gorm.DB

	count := 7
	dbUrl := getTestConnectionString()

	for i := 0; i < count; i++ {
		db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			if i == count-1 {
				os.Exit(FailedToConnect)
			}

			// failed to connect to db
			waitTime := time.Duration(i+1) * 5 * time.Second
			fmt.Printf("🛑 Failed to connect to database %s, attempting again after %s seconds\n", dbUrl, waitTime)
			time.Sleep(waitTime)
		} else {
			break
		}
	}

	cnn, _ := db.DB()
	err = cnn.Ping()
	if err != nil {
		fmt.Println("🛑 Failed to ping database")
		os.Exit(FailedToPing)
	}

	if err = AutoMigrate(db); err != nil {
		fmt.Println("🛑 Failed to run migrations")
		os.Exit(FailedToRunMigrations)
	}

	err = cnn.Close()
	if err != nil {
		fmt.Println("🛑 Failed to close connection")
		os.Exit(FailedToCloseConnection)
	}

	res := m.Run()

	os.Exit(res)
}
