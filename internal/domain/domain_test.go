package domain

import (
	"context"
	"errors"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/models/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignUp(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	user := types.User{
		ID:       0,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
		Payload:  nil,
	}

	toCreateUser := types.User{
		ID:       0,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
		Payload:  nil,
	}

	createdUser := types.User{
		ID:       1422,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "$2a$14$AR/iuSz0NDtdoY4Xv76KauYkqjAhBtUZg2klm38HfKBjaC8Qwe59m",
		Role:     "admin",
		Payload:  nil,
	}

	tests := []struct {
		name                      string
		mockGetByEmailError       error
		mockCreateUserArgUser     *types.User
		mockCreateUserReturnUser  types.User
		mockCreateUserReturnError error
		returnUser                types.User
		returnError               error
	}{
		{
			name: "user already exists",
		},
		{
			name:                "find user error",
			mockGetByEmailError: errors.New("some other error"),
		},
		{
			name:                      "create user error",
			mockGetByEmailError:       models.ErrNotFound,
			mockCreateUserArgUser:     &toCreateUser,
			mockCreateUserReturnUser:  types.User{},
			mockCreateUserReturnError: errors.New("create user error"),
		},
		{
			name:                     "success",
			mockGetByEmailError:      models.ErrNotFound,
			mockCreateUserArgUser:    &toCreateUser,
			mockCreateUserReturnUser: createdUser,
			returnUser:               createdUser,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockRepository := mocks.NewRepository(t)
			mockRepoUser := mocks.NewRepositoryUser(t)

			var userToReturn types.User
			if test.mockGetByEmailError != nil {
				userToReturn = types.User{}
			} else {
				userToReturn = user
			}

			userExpector := mockRepoUser.EXPECT()
			userExpector.GetByEmail(user.Email).Return(userToReturn, test.mockGetByEmailError)
			if test.mockCreateUserArgUser != nil {
				userExpector.Create(mock.Anything).
					Return(test.mockCreateUserReturnUser, test.mockCreateUserReturnError)
			}

			mockRepository.EXPECT().User(ctx).Return(mockRepoUser)

			// Test
			domain := NewDomain(mockRepository)
			returnUser, _ := domain.SignUp(ctx, user)

			// Assertions
			require.Equal(t, returnUser, test.returnUser)
		})
	}
}
