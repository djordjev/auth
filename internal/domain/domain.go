package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/models"
	"github.com/djordjev/auth/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name Domain
type Domain interface {
	SignUp(ctx context.Context, user types.User) (newUser types.User, err error)
	LogIn(ctx context.Context, user types.User) (exisingUser types.User, err error)
	LogOut(ctx context.Context, user types.User) (loggedOut bool, err error)
	Delete(ctx context.Context, user types.User) (deleted bool, err error)
	Verify(ctx context.Context, user types.User) (verified bool, err error)
}

func NewDomain(repository models.Repository, config utils.Config) Domain {
	return &domain{db: repository, config: config}
}

type domain struct {
	db     models.Repository
	config utils.Config
}

func (d *domain) Verify(ctx context.Context, user types.User) (verified bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *domain) SignUp(ctx context.Context, user types.User) (newUser types.User, err error) {
	userModel := d.db.User(ctx)

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

	d.db.Atomic(func(txRepo models.Repository) error {
		newUser, err = txRepo.User(ctx).Create(user)
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

		_, e = txRepo.VerifyAccount(ctx).Create(token.String(), newUser.ID)
		if e != nil {
			return e
		}

		return nil
	})

	return
}

func (d *domain) LogIn(ctx context.Context, user types.User) (exisingUser types.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *domain) LogOut(ctx context.Context, user types.User) (loggedOut bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *domain) Delete(ctx context.Context, user types.User) (deleted bool, err error) {
	//TODO implement me
	panic("implement me")
}
