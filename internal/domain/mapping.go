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
	return User{
		ID:       user.ID,
		Email:    user.Email,
		Username: *user.Username,
		Password: user.Password,
		Role:     user.Role,
		Verified: user.Verified,
	}
}
