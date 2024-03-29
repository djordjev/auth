package api

import (
	"fmt"
	"net/http"

	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
)

type SignUpRequest struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Email    string         `json:"email"`
	Role     string         `json:"role"`
	Payload  map[string]any `json:"payload"`
}

type SignUpResponse struct {
	ID       uint64         `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
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

	err = validateSignup(req)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)
	user, err := a.domain.SignUp(setup, signUpRequestToUser(req))
	if err == domain.ErrUserAlreadyExists {
		utils.LogError(logger, err)
		respondWithError(w, fmt.Sprintf("user with email %s already exists", req.Email), http.StatusBadRequest)
		return
	} else if err != nil {
		utils.LogError(logger, err)
		respondWithInternalError(w)
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
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func (a *jsonApi) postLogin(w http.ResponseWriter, r *http.Request) {
	var req LogInRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		respondWithError(w, "bad request", http.StatusBadRequest)
		return
	}

	err = validateLogin(req)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)
	user, session, err := a.domain.LogIn(setup, logInRequestToUser(req))
	if err == domain.ErrInvalidCredentials {
		respondWithError(w, "invalid credentials", http.StatusBadRequest)
		return
	} else if err != nil {
		respondWithError(w, "failed login attempt", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     a.cfg.SessionCookie,
		Value:    session,
		Path:     "/",
		MaxAge:   int(utils.SESSION_TTL.Milliseconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	mustWriteJSONResponse(w, userToLogInResponse(user))
}

type DeleteAccountRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteAccountResponse struct {
	Success bool `json:"success"`
}

func (a *jsonApi) deleteAccount(w http.ResponseWriter, r *http.Request) {
	var req DeleteAccountRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		respondWithBadRequest(w)
		return
	}

	err = validateDeleteAccount(req)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)

	deleted, err := a.domain.Delete(setup, deleteRequestToUser(req))
	if err == domain.ErrUserNotExist {
		respondWithError(w, "user does not exists", http.StatusBadRequest)
		return
	} else if err == domain.ErrInvalidCredentials {
		respondWithError(w, "authentication failed", http.StatusBadRequest)
		return
	} else if err != nil {
		respondWithInternalError(w)
		return
	}

	mustWriteJSONResponse(w, DeleteAccountResponse{Success: deleted})
}

type ForgetPasswordRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ForgetPasswordResponse struct {
	Success bool   `json:"success"`
	Email   string `json:"email"`
}

func (a *jsonApi) postForgetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: this should be authnticated request
	var req ForgetPasswordRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		respondWithBadRequest(w)
		return
	}

	err = validateForgetPassword(req)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)
	user, err := a.domain.ResetPasswordRequest(setup, forgetPasswordToUser(req))
	if err == domain.ErrUserNotExist {
		respondWithError(w, "user does not exist", http.StatusBadRequest)
		return
	}

	if user.Email != "" {
		mustWriteJSONResponse(w, ForgetPasswordResponse{Success: true, Email: user.Email})
		return
	}

	mustWriteJSONResponse(w, ForgetPasswordResponse{Success: false})
}

type SessionResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func (a *jsonApi) getSession(w http.ResponseWriter, r *http.Request) {
	var token string
	logger := utils.MustGetLogger(r)

	tokenCookie, err := r.Cookie(a.cfg.SessionCookie)
	if err == nil {
		token = tokenCookie.Value
	} else {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		respondWithBadRequest(w)
		return
	}
	setup := domain.NewSetup(r.Context(), logger)
	user, err := a.domain.Session(setup, token)
	if err == domain.ErrNoSession {
		respondWithBadRequest(w)
		return
	}

	response := SessionResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
		Verified: user.Verified,
	}

	mustWriteJSONResponse(w, response)
}

func (a *jsonApi) postLogout(w http.ResponseWriter, r *http.Request) {
	var token string
	logger := utils.MustGetLogger(r)

	tokenCookie, err := r.Cookie(a.cfg.SessionCookie)
	if err == nil {
		token = tokenCookie.Value
	} else {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		respondWithBadRequest(w)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)
	err = a.domain.Logout(setup, token)
	if err == domain.ErrNoSession {
		respondWithBadRequest(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
