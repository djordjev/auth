package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"

	"gorm.io/gorm"
)

type ForgetPassword struct {
	ModelWithDeletes
	UserID uint
	User   User
	Token  string `gorm:"unique,not null"`
}

type repositoryForgetPassword struct {
	ctx context.Context
	db  *gorm.DB
}

func (fp *repositoryForgetPassword) Create(token string, userId uint) (request domain.ForgetPassword, err error) {
	req := ForgetPassword{
		Token:  token,
		UserID: userId,
	}

	result := fp.db.Where("user_id = ?", userId).Delete(&ForgetPassword{})
	if result.Error != nil {
		err = fmt.Errorf("failed to delete previous password change requests for user %d %w", userId, result.Error)
		return
	}

	result = fp.db.Create(&req)
	if result.Error != nil {
		err = fmt.Errorf("failed to create password request for user %d %w", userId, result.Error)
		return
	}

	request.ID = req.ID
	request.Token = token
	request.UserID = userId

	return
}

func (fp *repositoryForgetPassword) Delete(token string) (request domain.ForgetPassword, err error) {
	forgetPasswordReq := ForgetPassword{}

	stored := fp.db.Where("token = ?", token).First(&forgetPasswordReq)
	if stored.Error == gorm.ErrRecordNotFound {
		err = fmt.Errorf("there's no reset request associated with token %s %w", token, stored.Error)
		return
	} else if stored.Error != nil {
		err = fmt.Errorf("failed to find token %s %w", token, stored.Error)
		return
	}

	result := fp.db.Delete(&ForgetPassword{}, forgetPasswordReq.ID)
	if result.Error != nil {
		err = fmt.Errorf("failed to delete password reset request for user %d %w", forgetPasswordReq.UserID, result.Error)
		return
	}

	request.ID = forgetPasswordReq.ID
	request.Token = token
	request.UserID = forgetPasswordReq.UserID

	return
}

func newRepositoryForgetPassword(ctx context.Context, db *gorm.DB) *repositoryForgetPassword {
	return &repositoryForgetPassword{ctx: ctx, db: db}
}
