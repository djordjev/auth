package domain

import (
	"context"
	"errors"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		err = errors.New("user already exists")
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return
	}

	user.Password = string(hash)

	newUser, err = userModel.Create(user)
	if err != nil {
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
