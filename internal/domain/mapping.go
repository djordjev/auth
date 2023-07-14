package domain

import (
	"github.com/djordjev/auth/internal/models"
)

func userToModel(user User) models.User {
	model := models.User{
		Email:    user.Email,
		Password: user.Password,
		Username: &user.Username,
		Role:     user.Role,
		Verified: user.Verified,
	}

	model.ID = user.ID

	return model
}

func modelToUser(user models.User) User {
	usr := User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Verified: user.Verified,
	}

	if user.Username != nil {
		usr.Username = *user.Username
	} else {
		usr.Username = ""
	}

	return usr
}
