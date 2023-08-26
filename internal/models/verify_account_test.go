package models

import (
	"context"
	"testing"

	"github.com/djordjev/auth/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestVerifyAccountCreate(t *testing.T) {
	existingUser := newRandomUser()
	token, _ := uuid.NewUUID()

	res := dbConnection.Create(&existingUser)
	require.Nil(t, res.Error, "failed to initialize db state")

	repo := newRepositoryVerifyAccount(context.TODO(), dbConnection)

	type testCase struct {
		name   string
		userId uint
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

func TestVerifyAccountVerify(t *testing.T) {
	existingUser := newRandomUser()
	token, _ := uuid.NewUUID()

	res := dbConnection.Create(&existingUser)
	require.Nil(t, res.Error, "failed to initialize db state")

	verification := VerifyAccount{UserID: existingUser.ID, Token: token.String()}
	res = dbConnection.Create(&verification)
	require.Nil(t, res.Error, "failed to initialize db state")

	repo := newRepositoryVerifyAccount(context.TODO(), dbConnection)

	type testCase struct {
		name        string
		token       string
		result      domain.VerifyAccount
		resultError string
	}

	tests := []testCase{
		{
			name:  "verifies",
			token: token.String(),
			result: domain.VerifyAccount{
				Token:  token.String(),
				UserID: existingUser.ID,
			},
		},
		{
			name:        "invalid token",
			token:       "abc",
			resultError: "no verify request associated with token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			result, err := repo.Verify(tc.token)

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
