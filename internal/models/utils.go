package models

import (
	"fmt"

	"github.com/djordjev/pg-mig/migrations"
)

func AutoMigrate(connectionString string) error {

	fmt.Println(connectionString)

	runner := migrations.NewRunner(
		"localhost",
		"tester:testee",
		"testdb",
		5433,
		"/Users/djordjev/Documents/Projects/auth/migrations",
	)

	err := runner.Run([]string{})

	return err
}
