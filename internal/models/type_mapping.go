package models

import (
	"github.com/djordjev/auth/internal/domain"
)

func modelUserToDomainUser(model User) domain.User {
	return domain.User{
		ID:       model.ID,
		Email:    model.Email,
		Username: *model.Username,
		Password: model.Password,
		Role:     model.Role,
		Verified: model.Verified,
	}
}

func domainUserToModelUser(user domain.User) User {
	usr := User{
		Email:    user.Email,
		Password: user.Password,
		Username: &user.Username,
		Role:     user.Role,
		Verified: user.Verified,
	}

	usr.ID = user.ID

	return usr
}
