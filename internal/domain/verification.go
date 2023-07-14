package domain

import (
	"fmt"

	"github.com/djordjev/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (d *domain) VerifyAccount(setup Setup, token string) (verified bool, err error) {

	err = d.db.Atomic(func(txRepo models.Repository) error {
		verifyAccount, err := txRepo.VerifyAccount(setup.ctx).Verify(token)
		if err == models.ErrNotFound {
			return ErrInvalidToken
		} else if err != nil {
			return fmt.Errorf("failed to verify account %w", err)
		}

		err = txRepo.User(setup.ctx).Verify(verifyAccount.User)
		if err != nil {
			return fmt.Errorf("failed to verify account %w", err)
		}

		return nil
	})

	verified = err != nil
	return
}

func (d *domain) VerifyPasswordReset(setup Setup, token string, password string) (updated User, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		err = fmt.Errorf("domain VerifyPasswordReset -> failed to generate password hash %w", err)
		return
	}

	err = d.db.Atomic(func(txRepo models.Repository) error {
		resetRequest, err := txRepo.ForgetPassword(setup.ctx).Delete(token)
		if err == models.ErrNotFound {
			return ErrInvalidToken
		} else if err != nil {
			return fmt.Errorf("failed to reset password %w", err)
		}

		user := resetRequest.User
		updateErr := txRepo.User(setup.ctx).SetPassword(user, string(hash))
		if updateErr != nil {
			return fmt.Errorf("failed to reset password %w", err)
		}

		updated = modelToUser(user)

		return nil
	})

	return
}
