package domain

import (
	"context"
	"errors"
	"testing"

	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errModel = errors.New("model error")

func TestSignUp(t *testing.T) {
	type testCase struct {
		name            string
		user            User
		created         User
		config          utils.Config
		setupUserRepo   func(*MockRepositoryUser, *testCase)
		setupVerifyRepo func(*MockRepositoryVerifyAccount, *testCase)
		setupRepo       func(*MockRepository, *MockRepositoryUser, *MockRepositoryVerifyAccount, *testCase)
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

	returnUser := signUpUser
	returnUser.ID = 55

	tests := []testCase{
		{
			name: "create user returns error",
			user: signUpUser,
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(User{}, modelErrors.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(User{}, errors.New("create failed"))
			},
			setupVerifyRepo: func(rva *MockRepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *MockRepository, u *MockRepositoryUser, v *MockRepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(Repository) error) error {
					return f(r)
				})
			},
			returnUser:  User{},
			returnError: errors.New("create failed"),
		},
		{
			name: "test user exists",
			user: signUpUser,
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(User{}, nil)
			},
			setupVerifyRepo: func(rva *MockRepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *MockRepository, u *MockRepositoryUser, v *MockRepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)
			},
			returnUser:  User{},
			returnError: ErrUserAlreadyExists,
		},
		{
			name:    "success",
			user:    signUpUser,
			created: returnUser,
			config:  utils.Config{RequireVerification: true},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(User{}, modelErrors.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(tc.created, nil)
			},
			setupVerifyRepo: func(rva *MockRepositoryVerifyAccount, tc *testCase) {
				rva.EXPECT().Create(mock.Anything, returnUser.ID).Return(VerifyAccount{}, nil)
			},
			setupRepo: func(r *MockRepository, u *MockRepositoryUser, v *MockRepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)
				r.EXPECT().VerifyAccount(setup.ctx).Return(v)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(Repository) error) error {
					return f(r)
				})
			},
			returnUser:  returnUser,
			returnError: nil,
		},
		{
			name:    "success without verification",
			user:    signUpUser,
			created: returnUser,
			config:  utils.Config{RequireVerification: false},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.user.Email).Return(User{}, modelErrors.ErrNotFound)
				ru.EXPECT().Create(mock.Anything).Return(tc.created, nil)
			},
			setupVerifyRepo: func(rva *MockRepositoryVerifyAccount, tc *testCase) {},
			setupRepo: func(r *MockRepository, u *MockRepositoryUser, v *MockRepositoryVerifyAccount, tc *testCase) {
				r.EXPECT().User(setup.ctx).Return(u)

				r.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(Repository) error) error {
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
			repository := NewMockRepository(t)
			userRepository := NewMockRepositoryUser(t)
			verifyRepository := NewMockRepositoryVerifyAccount(t)

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

func TestLogIn(t *testing.T) {
	t.Parallel()

	setup := Setup{ctx: context.TODO(), logger: utils.NewSilentLogger()}

	type testCase struct {
		name          string
		inputUser     User
		setupUserRepo func(*MockRepositoryUser, *testCase)
		returnUser    User
		returnError   error
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("testee"), 14)

	existing := User{
		ID:       452,
		Email:    "djvukovic@gmail.com",
		Password: string(hash),
		Role:     "admin",
		Username: "djvukovic",
		Verified: true,
	}

	tests := []testCase{
		{
			name:      "success",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(existing, nil)
			},
			returnUser: existing,
		},
		{
			name:      "username does not exist",
			inputUser: User{Username: "djvukovic", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(User{}, errModel)
			},
			returnUser:  User{},
			returnError: ErrUserNotExist,
		},
		{
			name:      "email does not exist",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(User{}, errModel)
			},
			returnUser:  User{},
			returnError: ErrUserNotExist,
		},
		{
			name:      "incorrect password",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				incorrectPassword := existing
				incorrectPassword.Password = "incorrect"

				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(incorrectPassword, nil)
			},
			returnUser:  User{},
			returnError: ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			repository := NewMockRepository(t)
			userRepository := NewMockRepositoryUser(t)

			// Setup mocks
			repository.EXPECT().User(context.TODO()).Return(userRepository).Maybe()
			tc.setupUserRepo(userRepository, &tc)

			// Run
			domain := NewDomain(repository, utils.Config{})
			user, err := domain.LogIn(setup, tc.inputUser)

			// Assertions
			if tc.returnError != nil {
				require.ErrorIs(t, err, tc.returnError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, user, tc.returnUser)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	setup := Setup{ctx: context.TODO(), logger: utils.NewSilentLogger()}
	hash, _ := bcrypt.GenerateFromPassword([]byte("testee"), 14)

	existing := User{
		ID:       452,
		Email:    "djvukovic@gmail.com",
		Password: string(hash),
		Role:     "admin",
		Username: "djvukovic",
		Verified: true,
	}

	type testCase struct {
		name          string
		inputUser     User
		setupUserRepo func(*MockRepositoryUser, *testCase)
		returnDeleted bool
		returnError   error
	}

	tests := []testCase{
		{
			name:      "success",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(existing, nil)
				ru.EXPECT().Delete(existing.ID).Return(true, nil)
			},
			returnDeleted: true,
		},
		{
			name:      "username not exist",
			inputUser: User{Username: "djvukovic", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(User{}, errModel)
			},
			returnError: ErrUserNotExist,
		},
		{
			name:      "email not exist",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(User{}, errModel)
			},
			returnError: ErrUserNotExist,
		},
		{
			name:      "delete failed",
			inputUser: User{Username: "djvukovic", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(existing, nil)
				ru.EXPECT().Delete(existing.ID).Return(false, errModel)
			},
			returnError: errModel,
		},
		{
			name:      "password incorrect",
			inputUser: User{Email: "djvukovic@gmail.com", Password: "testee"},
			setupUserRepo: func(ru *MockRepositoryUser, tc *testCase) {
				incorrectPassword := existing
				incorrectPassword.Password = "incorrect"

				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(incorrectPassword, nil)
			},
			returnError: ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			repository := NewMockRepository(t)
			userRepository := NewMockRepositoryUser(t)

			// Setup mocks
			repository.EXPECT().User(context.TODO()).Return(userRepository).Maybe()
			tc.setupUserRepo(userRepository, &tc)

			// Run
			domain := NewDomain(repository, utils.Config{})
			deleted, err := domain.Delete(setup, tc.inputUser)

			// Assertions
			require.Equal(t, deleted, tc.returnDeleted)
			if tc.returnError != nil {
				require.ErrorIs(t, err, tc.returnError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPasswordRequest(t *testing.T) {
	t.Parallel()

	setup := Setup{ctx: context.TODO(), logger: utils.NewSilentLogger()}
	hash, _ := bcrypt.GenerateFromPassword([]byte("testee"), 14)

	existing := User{

		ID:       452,
		Email:    "djvukovic@gmail.com",
		Password: string(hash),
		Role:     "admin",
		Username: "djvukovic",
		Verified: true,
	}

	type testCase struct {
		name        string
		inputUser   User
		setupModels func(*MockRepositoryUser, *MockRepositoryForgetPassword, *testCase)
		returnUser  User
		returnError error
	}

	tests := []testCase{
		{
			name:      "success",
			inputUser: User{Username: "djvukovic"},
			setupModels: func(ru *MockRepositoryUser, rfp *MockRepositoryForgetPassword, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(existing, nil)
				rfp.EXPECT().Create(mock.Anything, existing.ID).Return(ForgetPassword{}, nil)
			},
			returnUser: existing,
		},
		{
			name:      "username not exist",
			inputUser: User{Username: "djvukovic"},
			setupModels: func(ru *MockRepositoryUser, rfp *MockRepositoryForgetPassword, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(User{}, modelErrors.ErrNotFound)
			},
			returnError: ErrUserNotExist,
		},
		{
			name:      "email not exist",
			inputUser: User{Email: "djvukovic@gmail.com"},
			setupModels: func(ru *MockRepositoryUser, rfp *MockRepositoryForgetPassword, tc *testCase) {
				ru.EXPECT().GetByEmail(tc.inputUser.Email).Return(User{}, modelErrors.ErrNotFound)
			},
			returnError: ErrUserNotExist,
		},
		{
			name:      "failed to create reset request",
			inputUser: User{Username: "djvukovic"},
			setupModels: func(ru *MockRepositoryUser, rfp *MockRepositoryForgetPassword, tc *testCase) {
				ru.EXPECT().GetByUsername(tc.inputUser.Username).Return(existing, nil)
				rfp.EXPECT().Create(mock.Anything, existing.ID).Return(ForgetPassword{}, errModel)
			},
			returnError: errModel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			repository := NewMockRepository(t)
			userRepository := NewMockRepositoryUser(t)
			fpRepository := NewMockRepositoryForgetPassword(t)

			// Setup mocks
			repository.EXPECT().User(context.TODO()).Return(userRepository).Maybe()
			repository.EXPECT().ForgetPassword(context.TODO()).Return(fpRepository).Maybe()
			tc.setupModels(userRepository, fpRepository, &tc)

			// Run
			domain := NewDomain(repository, utils.Config{})
			user, err := domain.ResetPasswordRequest(setup, tc.inputUser)

			// Assertions
			require.Equal(t, user, tc.returnUser)
			if tc.returnError != nil {
				require.ErrorIs(t, err, tc.returnError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
