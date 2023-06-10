package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
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

type RepositoryUserTestSuite struct {
	suite.Suite
	conn *gorm.DB
}

func (suite *RepositoryUserTestSuite) SetupSuite() {
	conn := getTestConnectionString()
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		suite.FailNow("unable to acquire connection")
		return
	}

	suite.conn = db
}

func (suite *RepositoryUserTestSuite) TestGetByEmailExists() {
	existingUser := newRandomUser()

	res := suite.conn.Create(&existingUser)
	if res.Error != nil {
		suite.FailNow("failed to create user")
		return
	}

	repo := newRepositoryUser(context.TODO(), suite.conn)
	user, err := repo.GetByEmail(existingUser.Email)

	suite.Require().Nil(err)
	suite.Require().Equal(user.Role, existingUser.Role)
	suite.Require().Equal(user.Email, existingUser.Email)
	suite.Require().Equal(user.Password, existingUser.Password)
	suite.Require().Equal(user.Verified, existingUser.Verified)
	suite.Require().NotZero(user.ID)
}

func (suite *RepositoryUserTestSuite) TestGetByEmailNotExists() {
	repo := newRepositoryUser(context.TODO(), suite.conn)
	user, err := repo.GetByEmail(uuid.New().String())

	suite.Require().ErrorIs(err, ErrNotFound)
	suite.Require().Equal(user, User{})
}

func (suite *RepositoryUserTestSuite) TestGetByUsernameExists() {
	existingUser := newRandomUser()

	res := suite.conn.Create(&existingUser)
	if res.Error != nil {
		suite.FailNow("failed to create user")
		return
	}

	repo := newRepositoryUser(context.TODO(), suite.conn)
	user, err := repo.GetByUsername(*existingUser.Username)

	suite.Require().Nil(err)
	suite.Require().Equal(user.Role, existingUser.Role)
	suite.Require().Equal(user.Email, existingUser.Email)
	suite.Require().Equal(user.Password, existingUser.Password)
	suite.Require().Equal(user.Verified, existingUser.Verified)
	suite.Require().NotZero(user.ID)
}

func (suite *RepositoryUserTestSuite) TestGetByUsernameNotExists() {
	repo := newRepositoryUser(context.TODO(), suite.conn)
	user, err := repo.GetByUsername(uuid.New().String())

	suite.Require().ErrorIs(err, ErrNotFound)
	suite.Require().Equal(user, User{})
}

func (suite *RepositoryUserTestSuite) TestCreateSuccess() {
	repo := newRepositoryUser(context.TODO(), suite.conn)

	newUser := newRandomUser()

	user, err := repo.Create(newUser)

	suite.Require().Nil(err)
	suite.Require().Equal(user.Role, newUser.Role)
	suite.Require().Equal(user.Email, newUser.Email)
	suite.Require().Equal(user.Password, newUser.Password)
	suite.Require().Equal(user.Verified, newUser.Verified)
	suite.Require().NotZero(user.ID)
}

func (suite *RepositoryUserTestSuite) TestCreateError() {
	newUser := newRandomUser()
	res := suite.conn.Create(&newUser)
	if res.Error != nil {
		suite.FailNow("failed to create user")
		return
	}

	repo := newRepositoryUser(context.TODO(), suite.conn)
	user, err := repo.Create(newUser)

	suite.Require().Contains(err.Error(), "unable to create user")
	suite.Require().Equal(user, User{})

}

func TestRepositoryUserTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryUserTestSuite))
}
