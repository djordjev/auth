package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"
	"gorm.io/gorm"
)

type VerifyAccount struct {
	gorm.Model
	UserID uint
	User   User
	Token  string `gorm:"unique,not null"`
}

type repositoryVerifyAccount struct {
	ctx context.Context
	db  *gorm.DB
}

func (v *repositoryVerifyAccount) Create(token string, userId uint) (verification domain.VerifyAccount, err error) {
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

func (v *repositoryVerifyAccount) Verify(token string) (verification domain.VerifyAccount, err error) {
	//TODO implement me
	panic("implement me")
}

func newRepositoryVerifyAccount(ctx context.Context, db *gorm.DB) *repositoryVerifyAccount {
	return &repositoryVerifyAccount{ctx: ctx, db: db}
}
