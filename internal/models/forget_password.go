package models

import (
	"context"

	"gorm.io/gorm"
)

type ForgetPassword struct {
	gorm.Model
	UserID uint
	User   User
	Token  string `gorm:"unique,not null"`
}

//go:generate mockery --name RepositoryForgetPassword
type RepositoryForgetPassword interface {
	Create(userId uint) (request ForgetPassword, err error)
	Delete(token string) (request ForgetPassword, err error)
}

type repositoryForgetPassword struct {
	ctx context.Context
	db  *gorm.DB
}

func newRepositoryForgetPassword(ctx context.Context, db *gorm.DB) RepositoryForgetPassword {
	return &repositoryForgetPassword{ctx: ctx, db: db}
}

func (fp *repositoryForgetPassword) Create(userId uint) (request ForgetPassword, err error) {
	panic("not implemented") // TODO: Implement
}

func (fp *repositoryForgetPassword) Delete(token string) (request ForgetPassword, err error) {
	panic("not implemented") // TODO: Implement
}
