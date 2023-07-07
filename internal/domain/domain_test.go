package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/models/mocks"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	type testCase struct {
		name            string
		user            User
		created         models.User
		config          utils.Config
		setupUserRepo   func(*mocks.RepositoryUser, *testCase)
		setupVerifyRepo func(*mocks.RepositoryVerifyAccount, *testCase)
		setupRepo       func(*mocks.Repository, *mocks.RepositoryUser, *mocks.RepositoryVerifyAccount, *testCase)
		returnUser      User
		returnError     error
	}

	setup := Setup{ctx: context.TODO(), logger: utils.NewSilentLogger()}
	signUpUser := User{
		ID:       0,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}

	modelSignUpUser := userToModel(signUpUser)
	modelSignUpUser.ID = 55

	returnUser := signUpUser
	returnUser.ID = 55

	tests := []testCase{
		{
			name: "create user returns error",
			user: signUpUser,
			setupUserRepo: func(ru *mocks.RepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(models.User{}, models.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(models.User{}, errors.New("create failed"))
			},
			setupVerifyRepo: func(rva *mocks.RepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *mocks.Repository, u *mocks.RepositoryUser, v *mocks.RepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(models.Repository) error) error {
					return f(r)
				})
			},
			returnUser:  User{},
			returnError: errors.New("create failed"),
		},
		{
			name: "test user exists",
			user: signUpUser,
			setupUserRepo: func(ru *mocks.RepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(models.User{}, nil)
			},
			setupVerifyRepo: func(rva *mocks.RepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *mocks.Repository, u *mocks.RepositoryUser, v *mocks.RepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)
			},
			returnUser:  User{},
			returnError: ErrUserAlreadyExists,
		},
		{
			name:    "success",
			user:    signUpUser,
			created: modelSignUpUser,
			config:  utils.Config{RequireVerification: true},
			setupUserRepo: func(ru *mocks.RepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(models.User{}, models.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(tc.created, nil)
			},
			setupVerifyRepo: func(rva *mocks.RepositoryVerifyAccount, tc *testCase) {
				rva.EXPECT().Create(mock.Anything, modelSignUpUser.ID).Return(models.VerifyAccount{}, nil)
			},
			setupRepo: func(r *mocks.Repository, u *mocks.RepositoryUser, v *mocks.RepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)
				r.EXPECT().VerifyAccount(setup.ctx).Return(v)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(models.Repository) error) error {
					return f(r)
				})
			},
			returnUser:  returnUser,
			returnError: nil,
		},
		{
			name:    "success without verification",
			user:    signUpUser,
			created: modelSignUpUser,
			config:  utils.Config{RequireVerification: false},
			setupUserRepo: func(ru *mocks.RepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(models.User{}, models.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(tc.created, nil)
			},
			setupVerifyRepo: func(rva *mocks.RepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *mocks.Repository, u *mocks.RepositoryUser, v *mocks.RepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(models.Repository) error) error {
					return f(r)
				})
			},
			returnUser:  returnUser,
			returnError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			repository := mocks.NewRepository(t)
			userRepository := mocks.NewRepositoryUser(t)
			verifyRepository := mocks.NewRepositoryVerifyAccount(t)

			// Setup mocks
			tc.setupRepo(repository, userRepository, verifyRepository, &tc)
			tc.setupUserRepo(userRepository, &tc)
			tc.setupVerifyRepo(verifyRepository, &tc)

			// Run
			domain := NewDomain(repository, tc.config)
			usr, err := domain.SignUp(setup, tc.user)

			// Assertions
			if tc.returnError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, usr, tc.returnUser)

		})
	}
}
