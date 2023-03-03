package models

import (
	"context"
	"fmt"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func newRandomUser() User {
	email := fmt.Sprintf("djordje.vukovic+%s@gmail.com", uuid.New())
	username := fmt.Sprintf("random_username_%s", uuid.New())

	return User{
		Email:    email,
		Password: "password",
		Username: &username,
		Role:     "admin",
	}
}

func TestGetByEmail(t *testing.T) {
	// Prepare database
	conn := getTestConnectionString()
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	
	if err != nil {
		require.FailNow(t, "unable to acquire connection")
		return
	}

	existingUser := newRandomUser()
	res := db.Create(&existingUser)
	if res.Error != nil {
		require.FailNow(t, "failed to create user")
		return
	}

	// Tests
	tests := []struct {
		name      string
		email     string
		result    types.User
		resultErr error
	}{
		{
			name:  "exists in the database",
			email: existingUser.Email,
			result: types.User{
				ID:       0,
				Email:    existingUser.Email,
				Username: *existingUser.Username,
				Password: existingUser.Password,
				Role:     existingUser.Role,
			},
		},
		{
			name:      "does not exist in the database",
			email:     fmt.Sprintf("doesnotexist+%s@gmail.com", uuid.New()),
			result:    types.User{},
			resultErr: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := newRepositoryUser(context.TODO(), db)

			user, error := repo.GetByEmail(test.email)

			if test.resultErr != nil {
				require.ErrorIs(t, error, test.resultErr)
			} else {
				require.Equal(t, user.Email, test.result.Email)
				require.Equal(t, user.Username, test.result.Username)
				require.Equal(t, user.Password, test.result.Password)
				require.Equal(t, user.Role, test.result.Role)
				require.NotZero(t, user.ID)
			}

		})
	}

}
