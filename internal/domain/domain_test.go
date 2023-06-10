package domain

import (
	"context"
	"errors"
	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/models/mocks"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DomainSignUpTestSuite struct {
	suite.Suite
	repository       *mocks.Repository
	userRepository   *mocks.RepositoryUser
	verifyRepository *mocks.RepositoryVerifyAccount
	setup            Setup
	user             User
}

func (suite *DomainSignUpTestSuite) SetupTest() {
	suite.repository = mocks.NewRepository(suite.T())
	suite.userRepository = mocks.NewRepositoryUser(suite.T())
	suite.verifyRepository = mocks.NewRepositoryVerifyAccount(suite.T())
	suite.setup = Setup{
		ctx:    context.TODO(),
		logger: utils.NewSilentLogger(),
	}

	suite.user = User{
		ID:       0,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}
}

func (suite *DomainSignUpTestSuite) TestCreateUserSuccess() {
	modelUser := userToModel(suite.user)
	modelUser.Password = "XYZ"
	modelUser.ID = 54321

	suite.repository.EXPECT().User(suite.setup.ctx).Return(suite.userRepository)

	suite.userRepository.EXPECT().GetByEmail(suite.user.Email).Return(models.User{}, models.ErrNotFound)
	suite.userRepository.EXPECT().Create(mock.Anything).Return(modelUser, nil)

	suite.repository.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(models.Repository) error) error {
		f(suite.repository)
		return nil
	})

	cfg := utils.Config{RequireVerification: false}

	testDomain := NewDomain(suite.repository, cfg)

	usr, err := testDomain.SignUp(suite.setup, suite.user)

	suite.Require().Nil(err)
	suite.Require().Equal(usr, User{
		ID:       modelUser.ID,
		Email:    modelUser.Email,
		Username: *modelUser.Username,
		Password: modelUser.Password,
		Role:     modelUser.Role,
		Verified: false,
	})

}

func (suite *DomainSignUpTestSuite) TestUserExists() {
	suite.repository.EXPECT().User(suite.setup.ctx).Return(suite.userRepository)

	suite.userRepository.EXPECT().GetByEmail(suite.user.Email).Return(userToModel(suite.user), nil)

	testDomain := NewDomain(suite.repository, utils.Config{})
	usr, err := testDomain.SignUp(suite.setup, suite.user)

	suite.Require().Zero(usr)
	suite.Require().ErrorIs(err, ErrUserAlreadyExists)
}

func (suite *DomainSignUpTestSuite) TestError() {
	returnError := errors.New("database error")

	suite.repository.EXPECT().User(suite.setup.ctx).Return(suite.userRepository)

	suite.userRepository.EXPECT().GetByEmail(suite.user.Email).Return(models.User{}, models.ErrNotFound)
	suite.userRepository.EXPECT().Create(mock.Anything).Return(userToModel(suite.user), returnError)

	suite.repository.EXPECT().Atomic(mock.Anything).RunAndReturn(func(f func(models.Repository) error) error {
		return f(suite.repository)
	})

	testDomain := NewDomain(suite.repository, utils.Config{})
	_, err := testDomain.SignUp(suite.setup, suite.user)

	suite.Require().ErrorIs(err, returnError)
}

func TestDomainSignUpTestSuite(t *testing.T) {
	suite.Run(t, new(DomainSignUpTestSuite))
}
