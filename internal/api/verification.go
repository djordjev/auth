package api

import (
	"net/http"

	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/utils"
)

type VerifyAccountRequest struct {
	Token string
}

type VerifyAccountResponse struct {
	Verified bool
}

func (a *jsonApi) postVerifyAccount(w http.ResponseWriter, r *http.Request) {
	var req VerifyAccountRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		respondWithBadRequest(w)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)

	verified, err := a.domain.VerifyAccount(setup, req.Token)
	if err == domain.ErrInvalidToken {
		respondWithError(w, "invalid verification token", http.StatusBadRequest)
		return
	}

	mustWriteJSONResponse(w, VerifyAccountResponse{Verified: verified})
}

type VerifyPasswordResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type VerifyPasswordResetResponse struct {
	Success bool `json:"success"`
}

func (a *jsonApi) postVerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req VerifyPasswordResetRequest
	logger := utils.MustGetLogger(r)

	err := parseRequest(r, &req)
	if err != nil {
		respondWithBadRequest(w)
		return
	}

	err = validateVerifyPasswordResetRequest(req)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	setup := domain.NewSetup(r.Context(), logger)
	user, err := a.domain.VerifyPasswordReset(setup, req.Token, req.NewPassword)
	if err == domain.ErrInvalidToken {
		respondWithError(w, "invalid token", http.StatusBadRequest)
		return
	}

	response := VerifyPasswordResetResponse{}
	response.Success = user.Password != ""

	mustWriteJSONResponse(w, response)
}
