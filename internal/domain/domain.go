package domain

import (
	"context"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/models"
)

type Domain interface {
	SignUp(ctx context.Context, user types.User) (newUser types.User, err error)
	LogIn(ctx context.Context, user types.User) (exisingUser types.User, err error)
	LogOut(ctx context.Context, user types.User) (loggedOut bool, err error)
	Delete(ctx context.Context, user types.User) (deleted bool, err error)
}

func NewDomain(repository models.Repository) *domain {
	return &domain{db: repository}
}

type domain struct {
	db models.Repository
}

func (d *domain) SignUp(ctx context.Context, user types.User) (newUser types.User, err error) {
	userModel := d.db.User(ctx)
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
