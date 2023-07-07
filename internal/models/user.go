package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"unique;uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Username *string `gorm:"unique"`
	Role     string  `gorm:"default:regular"`
	Verified bool    `gorm:"default:false"`
}

//go:generate mockery --name RepositoryUser
type RepositoryUser interface {
	Create(user User) (newUser User, err error)
	Delete(id uint) (success bool, err error)
	GetByEmail(email string) (user User, err error)
	GetByUsername(username string) (user User, err error)
	Verify(user User) error
	SetPassword(user User, password string) error
}

type repositoryUser struct {
	ctx context.Context
	db  *gorm.DB
}

func (r *repositoryUser) GetByEmail(email string) (user User, err error) {
	q := r.db.Session(&gorm.Session{})

	result := q.Where("email = ?", email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		err = ErrNotFound
		return
	} else if result.Error != nil {
		err = fmt.Errorf("model GetByEmail -> find user by email %s, %w", email, result.Error)
		return
	}

	return
}

func (r *repositoryUser) GetByUsername(username string) (user User, err error) {
	q := r.db.Session(&gorm.Session{})

	result := q.Where("username = ?", username).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		err = ErrNotFound
		return
	} else if result.Error != nil {
		err = fmt.Errorf("model GetByUsername -> find user by username %s, %w", username, result.Error)
		return
	}

	return
}

func (r *repositoryUser) Create(user User) (newUser User, err error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		err = fmt.Errorf("model Create -> unable to create user %w", result.Error)
		return
	}

	if result.RowsAffected != 1 {
		err = fmt.Errorf("model Create -> multiple rows affected %d", result.RowsAffected)
		return
	}

	return user, nil
}

func (r *repositoryUser) Delete(id uint) (success bool, err error) {
	return true, nil
}

func (r *repositoryUser) Verify(user User) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	result := r.db.Model(&user).Updates(User{Verified: true})
	if result.Error != nil {
		return fmt.Errorf("failed to verify user with ID %d %w", user.ID, result.Error)
	}

	return nil
}

func (r *repositoryUser) SetPassword(user User, password string) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	result := r.db.Model(&user).Updates(User{Password: password})
	if result.Error != nil {
		return fmt.Errorf("failed to set password for user with ID %d %w", user.ID, result.Error)
	}

	return nil
}

func newRepositoryUser(ctx context.Context, db *gorm.DB) *repositoryUser {
	return &repositoryUser{ctx: ctx, db: db}
}
