package api

import (
	"fmt"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
	"net/http"
)

type SignUpRequest struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Email    string         `json:"email"`
	Role     string         `json:"role"`
	Payload  map[string]any `json:"payload"`
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

	user, err := a.domain.SignUp(r.Context(), signUpRequestToUser(req))
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
