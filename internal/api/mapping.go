package api

import (
	"github.com/djordjev/auth/internal/domain/types"
)

func signUpRequestToUser(req SignUpRequest) types.User {
	return types.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		Payload:  req.Payload,
	}
}

func userToSignUpResponse(user types.User) SignUpResponse {
	return SignUpResponse{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		Payload:  user.Payload,
	}
}
