package domain

import (
	"context"
	"errors"
	"fmt"

	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"github.com/djordjev/auth/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Setup struct {
	ctx    context.Context
	logger *slog.Logger
}

func NewSetup(ctx context.Context, logger *slog.Logger) Setup {
	return Setup{
		ctx:    ctx,
		logger: logger,
	}
}

//go:generate mockery --name Domain
type Domain interface {
	SignUp(setup Setup, user User) (newUser User, err error)
	LogIn(setup Setup, user User) (existing User, err error)
	Delete(setup Setup, user User) (deleted bool, err error)
	VerifyAccount(setup Setup, token string) (verified bool, err error)
	ResetPasswordRequest(setup Setup, user User) (sentTo User, err error)
	VerifyPasswordReset(setup Setup, token string, password string) (updated User, err error)
}

func NewDomain(repository Repository, config utils.Config) Domain {
	return &domain{db: repository, config: config}
}

type domain struct {
	db     Repository
	config utils.Config
}

func (d *domain) LogIn(setup Setup, user User) (existingUser User, err error) {
	userModel := d.db.User(setup.ctx)

	if user.Username != "" {
		existingUser, err = userModel.GetByUsername(user.Username)
	} else {
		existingUser, err = userModel.GetByEmail(user.Email)
	}

	if err != nil {
		err = ErrUserNotExist
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		existingUser = User{}
		err = ErrInvalidCredentials
		return
	}

	return

}

func (d *domain) Delete(setup Setup, user User) (deleted bool, err error) {
	userModel := d.db.User(setup.ctx)

	var existingUser User
	if user.Username != "" {
		existingUser, err = userModel.GetByUsername(user.Username)
	} else {
		existingUser, err = userModel.GetByEmail(user.Email)
	}

	if err != nil {
		err = ErrUserNotExist
		deleted = false
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		deleted = false
		err = ErrInvalidCredentials
		return
	}

	deleted, err = userModel.Delete(existingUser.ID)
	return
}

func (d *domain) ResetPasswordRequest(setup Setup, user User) (sentTo User, err error) {
	forgetPasswordModel := d.db.ForgetPassword(setup.ctx)
	userModel := d.db.User(setup.ctx)

	if user.Email != "" {
		sentTo, err = userModel.GetByEmail(user.Email)
	} else {
		sentTo, err = userModel.GetByUsername(user.Username)
	}

	if err == modelErrors.ErrNotFound {
		err = ErrUserNotExist
		return
	} else if err != nil {
		err = fmt.Errorf("domain ResetPasswordRequest -> failed to fetch user")
		return
	}

	_, err = forgetPasswordModel.Create(sentTo.ID)

	if err != nil {
		return
	}

	return
}

func (d *domain) SignUp(setup Setup, user User) (newUser User, err error) {
	userModel := d.db.User(setup.ctx)

	_, err = userModel.GetByEmail(user.Email)
	if err == nil {
		err = ErrUserAlreadyExists
		return
	}

	if !errors.Is(err, modelErrors.ErrNotFound) {
		err = fmt.Errorf("domain SignUp -> error looking for user %s: %w", user.Email, err)
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		err = fmt.Errorf("domain SignUp -> failed to generate password hash %w", err)
		return
	}

	user.Password = string(hash)
	user.Verified = !d.config.RequireVerification

	err = d.db.Atomic(func(txRepo Repository) error {
		newUser, err := txRepo.User(setup.ctx).Create(user)
		if err != nil {
			err = fmt.Errorf("domain SignUp -> failed to create a new user %w", err)
			return err
		}

		if !d.config.RequireVerification {
			return nil
		}

		token, e := uuid.NewUUID()
		if e != nil {
			return e
		}

		_, e = txRepo.VerifyAccount(setup.ctx).Create(token.String(), newUser.ID)
		if e != nil {
			return e
		}

		return nil
	})

	return
}
