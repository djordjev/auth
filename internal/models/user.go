package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"
	modelErrors "github.com/djordjev/auth/internal/models/errors"

	"gorm.io/gorm"
)

type User struct {
	ModelWithDeletes
	Email    string  `gorm:"unique;uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Username *string `gorm:"unique;uniqueIndex"`
	Role     string  `gorm:"default:regular"`
	Verified bool    `gorm:"default:false"`
}

type repositoryUser struct {
	ctx context.Context
	db  *gorm.DB
}

func (r *repositoryUser) GetByEmail(email string) (user domain.User, err error) {
	q := r.db.Session(&gorm.Session{})

	modelUser := User{}
	result := q.Where("email = ?", email).First(&modelUser)

	if result.Error == gorm.ErrRecordNotFound {
		err = modelErrors.ErrNotFound
		return
	} else if result.Error != nil {
		err = fmt.Errorf("model GetByEmail -> find user by email %s, %w", email, result.Error)
		return
	}

	user = modelUserToDomainUser(modelUser)

	return
}

func (r *repositoryUser) GetByUsername(username string) (user domain.User, err error) {
	q := r.db.Session(&gorm.Session{})

	modelUser := User{}
	result := q.Where("username = ?", username).First(&modelUser)

	if result.Error == gorm.ErrRecordNotFound {
		err = modelErrors.ErrNotFound
		return
	} else if result.Error != nil {
		err = fmt.Errorf("model GetByUsername -> find user by username %s, %w", username, result.Error)
		return
	}

	user = modelUserToDomainUser(modelUser)

	return
}

func (r *repositoryUser) Create(user domain.User) (newUser domain.User, err error) {
	modelUser := domainUserToModelUser(user)

	result := r.db.Create(&modelUser)
	if result.Error != nil {
		err = fmt.Errorf("model Create -> unable to create user %w", result.Error)
		return
	}

	if result.RowsAffected != 1 {
		err = fmt.Errorf("model Create -> multiple rows affected %d", result.RowsAffected)
		return
	}

	newUser = modelUserToDomainUser(modelUser)

	return
}

func (r *repositoryUser) Delete(id uint) (success bool, err error) {
	user := User{}
	user.ID = id

	result := r.db.Delete(&user)

	if result.Error != nil {
		err = fmt.Errorf("model Delete -> failed to delete user with id %d %w", id, result.Error)
		return
	}

	return true, nil
}

func (r *repositoryUser) Verify(user domain.User) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	modelUser := domainUserToModelUser(user)
	result := r.db.Model(&modelUser).Updates(User{Verified: true})
	if result.Error != nil {
		return fmt.Errorf("failed to verify user with ID %d %w", user.ID, result.Error)
	}

	return nil
}

func (r *repositoryUser) SetPassword(user domain.User, password string) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	modelUser := domainUserToModelUser(user)
	result := r.db.Model(&modelUser).Updates(User{Password: password})
	if result.Error != nil {
		return fmt.Errorf("failed to set password for user with ID %d %w", user.ID, result.Error)
	}

	return nil
}

func newRepositoryUser(ctx context.Context, db *gorm.DB) *repositoryUser {
	return &repositoryUser{ctx: ctx, db: db}
}
