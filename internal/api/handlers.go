package api

import (
	"encoding/json"
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.domain.SignUp(r.Context(), signUpRequestToUser(req))
	if err != nil {
		http.Error(w, "error message", http.StatusInternalServerError)
		return
	}

	response := userToSignUpResponse(user)

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(responseBytes)
}
