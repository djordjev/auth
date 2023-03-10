package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name Domain
type Domain interface {
	SignUp(ctx context.Context, user types.User) (newUser types.User, err error)
	LogIn(ctx context.Context, user types.User) (exisingUser types.User, err error)
	LogOut(ctx context.Context, user types.User) (loggedOut bool, err error)
	Delete(ctx context.Context, user types.User) (deleted bool, err error)
}

func NewDomain(repository models.Repository) Domain {
	return &domain{db: repository}
}

type domain struct {
	db models.Repository
}

func (d *domain) SignUp(ctx context.Context, user types.User) (newUser types.User, err error) {
	userModel := d.db.User(ctx)

	_, err = userModel.GetByEmail(user.Email)
	if err == nil {
		userMessage := fmt.Sprintf("user with email %s already exists", user.Email)
		err = NewDomainError(userMessage, fmt.Sprintf("SignUp -> %s", userMessage), nil, false, user)
		return
	}

	if !errors.Is(err, models.ErrNotFound) {
		err = NewDomainError("", "SignUp -> error looking for the user", err, true, user)
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		err = NewDomainError("", "SignUp -> Unable to generate hash", err, true, user)
		return
	}

	user.Password = string(hash)

	newUser, err = userModel.Create(user)
	if err != nil {
		err = NewDomainError("", "SignUp -> unable to create user", err, true, user)
		return
	}

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
