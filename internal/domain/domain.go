package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/djordjev/auth/internal/models"
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
	LogIn(setup Setup, user User) (exisingUser User, err error)
	Delete(setup Setup, user User) (deleted bool, err error)
	VerifyAccount(setup Setup, token string) (verified bool, err error)
	ResetPasswordRequest(setup Setup, user User) (sentTo User, err error)
	VerifyPasswordReset(setup Setup, token string, password string) (updated User, err error)
}

func NewDomain(repository models.Repository, config utils.Config) Domain {
	return &domain{db: repository, config: config}
}

type domain struct {
	db     models.Repository
	config utils.Config
}

func (d *domain) LogIn(setup Setup, user User) (exisingUser User, err error) {
	userModel := d.db.User(setup.ctx)

	var modelUser models.User
	if user.Username != "" {
		modelUser, err = userModel.GetByUsername(user.Username)
	} else {
		modelUser, err = userModel.GetByEmail(user.Email)
	}

	if err != nil {
		err = fmt.Errorf("domain LogIn -> failed to find user with email %s and username %s, %w",
			user.Email,
			user.Username,
			err)

		return
	}

	exisingUser = modelToUser(modelUser)
	err = bcrypt.CompareHashAndPassword([]byte(exisingUser.Password), []byte(user.Password))
	if err != nil {
		err = ErrInvalidCredentials
		return
	}

	return

}

func (d *domain) Delete(setup Setup, user User) (deleted bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *domain) ResetPasswordRequest(setup Setup, user User) (sentTo User, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *domain) SignUp(setup Setup, user User) (newUser User, err error) {
	userModel := d.db.User(setup.ctx)

	_, err = userModel.GetByEmail(user.Email)
	if err == nil {
		err = ErrUserAlreadyExists
		return
	}

	if !errors.Is(err, models.ErrNotFound) {
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

	err = d.db.Atomic(func(txRepo models.Repository) error {
		modelUser, err := txRepo.User(setup.ctx).Create(userToModel(user))
		newUser = modelToUser(modelUser)
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
