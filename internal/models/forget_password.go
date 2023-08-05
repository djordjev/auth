package models

import (
	"context"
	"github.com/djordjev/auth/internal/domain"

	"gorm.io/gorm"
)

type ForgetPassword struct {
	gorm.Model
	UserID uint
	User   User
	Token  string `gorm:"unique,not null"`
}

type repositoryForgetPassword struct {
	ctx context.Context
	db  *gorm.DB
}

func (fp *repositoryForgetPassword) Create(userId uint) (request domain.ForgetPassword, err error) {
	//TODO implement me
	panic("implement me")
}

func (fp *repositoryForgetPassword) Delete(token string) (request domain.ForgetPassword, err error) {
	//TODO implement me
	panic("implement me")
}

func newRepositoryForgetPassword(ctx context.Context, db *gorm.DB) *repositoryForgetPassword {
	return &repositoryForgetPassword{ctx: ctx, db: db}
}
