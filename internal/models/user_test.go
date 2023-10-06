package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/djordjev/auth/internal/domain"
	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var nonExistingUserID = uint64(9223372036)

func newRandomUser() domain.User {
	email := fmt.Sprintf("djordje.vukovic+%s@gmail.com", uuid.New())
	username := fmt.Sprintf("random_username_%s", uuid.New())

	return domain.User{
		Email:    email,
		Password: "password",
		Username: username,
		Role:     "admin",
	}
}

func TestGetByEmail(t *testing.T) {
	existingUser := newRandomUser()

	existingUser, err := storeUser(existingUser)
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		email       string
		result      domain.User
		resultError error
	}

	tests := []testCase{
		{
			name:   "user found",
			email:  existingUser.Email,
			result: existingUser,
		},
		{
			name:        "user not found",
			email:       fmt.Sprintf("not_exist%s", existingUser.Email),
			resultError: modelErrors.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.GetByEmail(tc.email)

			require.Equal(t, result, tc.result)
			if tc.resultError != nil {
				require.ErrorIs(t, tc.resultError, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	existingUser := newRandomUser()

	existingUser, err := storeUser(existingUser)
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		username    string
		result      domain.User
		resultError error
	}

	tests := []testCase{
		{
			name:     "user found",
			username: existingUser.Username,
			result:   existingUser,
		},
		{
			name:        "user not found",
			username:    fmt.Sprintf("not_exist%s", existingUser.Email),
			resultError: modelErrors.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.GetByUsername(tc.username)

			require.Equal(t, result, tc.result)
			if tc.resultError != nil {
				require.ErrorIs(t, tc.resultError, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	newUser := newRandomUser()

	existingUser, err := storeUser(newRandomUser())
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		user        domain.User
		result      domain.User
		resultError string
	}

	tests := []testCase{
		{
			name:        "user already exists",
			user:        existingUser,
			resultError: "unable to create user",
		},
		{
			name:   "success",
			user:   newUser,
			result: newUser,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.Create(tc.user)

			require.Equal(t, result.Email, tc.result.Email)
			require.Equal(t, result.Username, tc.result.Username)
			require.Equal(t, result.Role, tc.result.Role)
			require.Equal(t, result.Verified, tc.result.Verified)

			if tc.resultError != "" {
				require.ErrorContains(t, err, tc.resultError)
			} else {
				require.Nil(t, err)
				require.NotZero(t, result.ID)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	existingUser, err := storeUser(newRandomUser())
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		id          uint64
		result      bool
		resultError string
	}

	tests := []testCase{
		{
			name:   "deletes existing user",
			id:     existingUser.ID,
			result: true,
		},
		{
			name:        "user does not exists",
			id:          nonExistingUserID,
			result:      false,
			resultError: "does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.Delete(tc.id)

			if tc.resultError != "" {
				require.ErrorContains(t, err, tc.resultError)
				require.False(t, result)
			} else {
				require.Nil(t, err)
				require.True(t, result)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	existingUser, err := storeUser(newRandomUser())
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name   string
		user   domain.User
		result string
	}

	tests := []testCase{
		{
			name:   "verifies user",
			user:   existingUser,
			result: "",
		},
		{
			name:   "user doesn't have ID",
			user:   domain.User{},
			result: "missing user ID in update function",
		},
		{
			name:   "fails to verify",
			user:   domain.User{ID: nonExistingUserID},
			result: "does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result := repo.Verify(tc.user)

			if tc.result != "" {
				require.ErrorContains(t, result, tc.result)
			} else {
				require.Nil(t, result)
			}
		})
	}
}

func TestSetPassword(t *testing.T) {
	existingUser, err := storeUser(newRandomUser())
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryUser(context.TODO(), dbConnection)

	type testCase struct {
		name   string
		user   domain.User
		result string
	}

	tests := []testCase{
		{
			name:   "updates password",
			user:   existingUser,
			result: "",
		},
		{
			name:   "user doesn't have ID",
			user:   domain.User{},
			result: "missing user ID in update function",
		},
		{
			name:   "fails to update",
			user:   domain.User{ID: nonExistingUserID},
			result: "does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result := repo.SetPassword(tc.user, "random_password")

			if tc.result != "" {
				require.ErrorContains(t, result, tc.result)
			} else {
				require.Nil(t, result)
			}
		})
	}
}
