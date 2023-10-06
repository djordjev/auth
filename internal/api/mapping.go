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
		Payload:  req.Payload,
	}
}

func userToSignUpResponse(user domain.User) SignUpResponse {
	return SignUpResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Payload:  user.Payload,
	}
}

func logInRequestToUser(req LogInRequest) domain.User {
	return domain.User{Email: req.Email, Username: req.Username, Password: req.Password}
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

func deleteRequestToUser(req DeleteAccountRequest) domain.User {
	return domain.User{Email: req.Email, Password: req.Password, Username: req.Username}
}

func forgetPasswordToUser(req ForgetPasswordRequest) domain.User {
	return domain.User{Email: req.Email, Username: req.Username}
}
