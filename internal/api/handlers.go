package api

import (
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

	err := parseRequest(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.domain.SignUp(r.Context(), signUpRequestToUser(req))
	if err != nil {
		respondWithError(w, err)

		logger := utils.MustGetLoggerFromRequest(r)
		logger.Info("something else")
		return
	}

	response := userToSignUpResponse(user)

	mustWriteJSONResponse(w, response)
}
