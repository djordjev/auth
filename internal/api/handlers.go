package api

import (
	"fmt"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
	"net/http"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type SignUpResponse struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Role     string         `json:"role"`
	Payload  map[string]any `json:"payload"`
}

func (a *jsonApi) postSignup(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.domain.SignUp(domain.NewSetup(r.Context(), logger), signUpRequestToUser(req))
	if err == domain.ErrUserAlreadyExists {
		utils.LogError(logger, err)
		http.Error(w, fmt.Sprintf("user with email %s already exists", req.Email), http.StatusBadRequest)
		return
	} else if err != nil {
		utils.LogError(logger, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := userToSignUpResponse(user)

	mustWriteJSONResponse(w, response)
}

type LogInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LogInResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"'`
}

func (a *jsonApi) postLogin(w http.ResponseWriter, r *http.Request) {
	var req LogInRequest
	err := parseRequest(r, &req)
	logger := utils.MustGetLogger(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.domain.LogIn(domain.NewSetup(r.Context(), logger), LogInRequestToUser(req))
	if err == domain.ErrInvalidCredentials {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "failed login attempt", http.StatusInternalServerError)
		return
	}

	mustWriteJSONResponse(w, userToLogInResponse(user))
}

type DeleteAccountRequest struct {
}

type DeleteAccountResponse struct {
}

func (a *jsonApi) deleteAccount(w http.ResponseWriter, r *http.Request) {}

type UpdateAccountRequest struct {
}

type UpdateAccountResponse struct {
}

func (a *jsonApi) putAccount(w http.ResponseWriter, r *http.Request) {}
