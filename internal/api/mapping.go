package api

import (
	"github.com/djordjev/auth/internal/domain"
)

func signUpRequestToUser(req SignUpRequest) domain.User {
	return domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}
}

func userToSignUpResponse(user domain.User) SignUpResponse {
	return SignUpResponse{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}

func LogInRequestToUser(req LogInRequest) domain.User {
	return domain.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}
}

func userToLogInResponse(user domain.User) LogInResponse {
	return LogInResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
		Verified: user.Verified,
	}
}
