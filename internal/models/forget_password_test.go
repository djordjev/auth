package models

import (
	"context"
	"testing"

	"github.com/djordjev/auth/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestForgetPasswordCreate(t *testing.T) {
	existingUser := newRandomUser()
	token, _ := uuid.NewUUID()

	existingUser, err := storeUser(existingUser)
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryForgetPassword(context.TODO(), dbConnection)

	type testCase struct {
		name   string
		userId uint64
		result domain.VerifyAccount
	}

	tests := []testCase{
		{
			name:   "creates new verification record",
			userId: existingUser.ID,
			result: domain.VerifyAccount{
				Token:  token.String(),
				UserID: existingUser.ID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.Create(token.String(), existingUser.ID)

			require.Nil(t, err)
			require.Equal(t, result.Token, token.String())
			require.Equal(t, result.UserID, existingUser.ID)
		})
	}
}

func TestForgetPasswordDelete(t *testing.T) {
	existingUser := newRandomUser()
	token, _ := uuid.NewUUID()

	existingUser, err := storeUser(existingUser)
	require.Nil(t, err, "failed to initialize db state")

	_, err = dbConnection.Exec(
		context.Background(),
		"insert into forget_passwords (user_id, token) values ($1, $2)",
		existingUser.ID, token.String(),
	)
	require.Nil(t, err, "failed to initialize db state")

	repo := newRepositoryForgetPassword(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		token       string
		result      domain.ForgetPassword
		resultError string
	}

	tests := []testCase{
		{
			name:  "verifies",
			token: token.String(),
			result: domain.ForgetPassword{
				Token:  token.String(),
				UserID: existingUser.ID,
			},
		},
		{
			name:        "invalid token",
			token:       "abc",
			resultError: "no reset request associated with token abc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.Delete(tc.token)

			require.Equal(t, result.Token, tc.result.Token)
			require.Equal(t, result.UserID, tc.result.UserID)

			if tc.resultError != "" {
				require.ErrorContains(t, err, tc.resultError)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
