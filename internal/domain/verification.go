package domain

import (
	"fmt"

	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"golang.org/x/crypto/bcrypt"
)

func (d *domain) VerifyAccount(setup Setup, token string) (verified bool, err error) {
	err = d.db.Atomic(func(txRepo Repository) error {
		verifyAccount, err := txRepo.VerifyAccount(setup.ctx).Verify(token)
		if err == modelErrors.ErrNotFound {
			return ErrInvalidToken
		} else if err != nil {
			return fmt.Errorf("failed to verify account %w", err)
		}

		user := User{ID: verifyAccount.ID}
		err = txRepo.User(setup.ctx).Verify(user)
		if err != nil {
			return fmt.Errorf("failed to verify account %w", err)
		}

		return nil
	})

	verified = err == nil
	return
}

func (d *domain) VerifyPasswordReset(setup Setup, token string, password string) (updated User, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		err = fmt.Errorf("domain VerifyPasswordReset -> failed to generate password hash %w", err)
		return
	}

	err = d.db.Atomic(func(txRepo Repository) error {
		resetRequest, err := txRepo.ForgetPassword(setup.ctx).Delete(token)
		if err == modelErrors.ErrNotFound {
			return ErrInvalidToken
		} else if err != nil {
			return fmt.Errorf("failed to reset password %w", err)
		}

		user := User{ID: resetRequest.UserID}
		newPassword := string(hash)
		updateErr := txRepo.User(setup.ctx).SetPassword(user, newPassword)
		if updateErr != nil {
			return fmt.Errorf("failed to reset password %w", err)
		}

		updated.Password = newPassword
		return nil
	})

	return
}
