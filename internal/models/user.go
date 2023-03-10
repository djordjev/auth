package models

import (
	"context"
	"fmt"
	"github.com/djordjev/auth/internal/domain/types"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"unique;uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Username *string `gorm:"unique"`
	Role     string  `gorm:"default:regular"`
}

//go:generate mockery --name RepositoryUser
type RepositoryUser interface {
	GetByEmail(email string) (user types.User, err error)
	Create(user types.User) (newUser types.User, err error)
}

type repositoryUser struct {
	ctx context.Context
	db  *gorm.DB
}

func (r *repositoryUser) GetByEmail(email string) (user types.User, err error) {
	q := r.db.Session(&gorm.Session{})

	dbUser := User{}
	result := q.Where("email = ?", email).First(&dbUser)

	if result.Error == gorm.ErrRecordNotFound {
		err = ErrNotFound
		return
	} else if result.Error != nil {
		err = NewModelError("GetByEmail -> find user by email", result.Error)
		return
	}

	user.ID = dbUser.ID
	user.Username = *dbUser.Username
	user.Password = dbUser.Password
	user.Role = dbUser.Role
	user.Email = dbUser.Email

	return
}

func (r *repositoryUser) Create(user types.User) (newUser types.User, err error) {
	usr := User{
		Email:    user.Email,
		Password: user.Password,
		Username: &user.Username,
		Role:     user.Role,
	}

	result := r.db.Create(&usr)
	if result.Error != nil {
		err = NewModelError("Create -> unable to create user", result.Error)
		return
	}

	if result.RowsAffected != 1 {
		err = NewModelError(fmt.Sprintf("Create -> multiple rows affected (%d)", result.RowsAffected), nil)
		return
	}

	newUser.ID = usr.ID
	newUser.Email = usr.Email
	if usr.Username != nil {
		newUser.Username = *usr.Username
	} else {
		newUser.Username = ""
	}
	newUser.Password = usr.Password
	newUser.Role = usr.Role

	return
}

func newRepositoryUser(ctx context.Context, db *gorm.DB) *repositoryUser {
	return &repositoryUser{ctx: ctx, db: db}
}
