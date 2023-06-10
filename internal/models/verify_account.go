package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type VerifyAccount struct {
	gorm.Model
	UserID uint
	User   User
	Token  string `gorm:"unique,not null"`
}

//go:generate mockery --name RepositoryVerifyAccount
type RepositoryVerifyAccount interface {
	Create(token string, userId uint) (verification VerifyAccount, err error)
	Verify(token string) (verification VerifyAccount, err error)
}

type repositoryVerifyAccount struct {
	ctx context.Context
	db  *gorm.DB
}

func (v *repositoryVerifyAccount) Create(token string, userId uint) (verification VerifyAccount, err error) {
	ver := VerifyAccount{
		UserID: userId,
		Token:  token,
	}

	result := v.db.Create(&ver)
	if result.Error != nil {
		err = fmt.Errorf("model VerifyAccount -> unable to create verification token")
		return
	}

	verification.Token = token
	verification.ID = ver.ID
	verification.UserID = userId

	return
}

func (v *repositoryVerifyAccount) Verify(token string) (verification VerifyAccount, err error) {
	//TODO implement me
	panic("implement me")
}

func newRepositoryVerifyAccount(ctx context.Context, db *gorm.DB) RepositoryVerifyAccount {
	return &repositoryVerifyAccount{ctx: ctx, db: db}
}
