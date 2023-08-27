package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"
	"gorm.io/gorm"
)

type VerifyAccount struct {
	ModelWithDeletes
	UserID uint
	User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	verifyRequest := VerifyAccount{}

	stored := v.db.Where("token = ?", token).First(&verifyRequest)
	if stored.Error == gorm.ErrRecordNotFound {
		err = fmt.Errorf("there's no verify request associated with token %s %w", token, stored.Error)
		return
	} else if stored.Error != nil {
		err = fmt.Errorf("failed to find verify token %s %w", token, stored.Error)
		return
	}

	result := v.db.Delete(&VerifyAccount{}, verifyRequest.ID)
	if result.Error != nil {
		err = fmt.Errorf("failed to delete verification request for user %d %w", verifyRequest.UserID, result.Error)
		return
	}

	verification.ID = verifyRequest.ID
	verification.Token = token
	verification.UserID = verifyRequest.UserID

	return
}

func newRepositoryVerifyAccount(ctx context.Context, db *gorm.DB) *repositoryVerifyAccount {
	return &repositoryVerifyAccount{ctx: ctx, db: db}
}
