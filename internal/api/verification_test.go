package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/domain/mocks"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var logger = utils.NewSilentLogger()

func TestVerifyAccount(t *testing.T) {
	var verifyAccountBody = `
		{ "token": "abc" }
	`

	tests := []struct {
		name       string
		request    string
		statusCode int
		response   string
		returnVal  bool
		returnErr  error
	}{
		{
			name:       "success",
			request:    verifyAccountBody,
			statusCode: http.StatusOK,
			response:   `{ "verified": true }`,
			returnVal:  true,
		},
		{
			name:       "wrong token",
			request:    verifyAccountBody,
			statusCode: http.StatusBadRequest,
			response:   utils.ErrorJSON("invalid verification token"),
			returnErr:  domain.ErrInvalidToken,
		},
		{
			name:       "internal error",
			request:    verifyAccountBody,
			statusCode: http.StatusInternalServerError,
			response:   utils.ErrorJSON("internal server error"),
			returnErr:  errors.New("random error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := utils.RequestBuilder("POST", "/verify")(tc.request)

			rr := httptest.NewRecorder()

			baseExpector := mocks.NewDomain(t)
			domain := baseExpector.EXPECT()

			matcher := mock.MatchedBy(func(token string) bool {
				return token == "abc"
			})

			domain.VerifyAccount(mock.Anything, matcher).Return(tc.returnVal, tc.returnErr)

			api := NewApi(utils.Config{}, http.NewServeMux(), baseExpector, logger)

			api.postVerifyAccount(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)
			require.JSONEq(t, tc.response, rr.Body.String())
		})
	}
}

func TestVerifyPasswordReset(t *testing.T) {
	var verifyPasswordReset = `
	{ "token": "abc", "new_password": "testee" }
`

	tests := []struct {
		name       string
		request    string
		statusCode int
		response   string
		returnVal  domain.User
		returnErr  error
		skipMock   bool
	}{
		{
			name:       "success",
			request:    verifyPasswordReset,
			statusCode: http.StatusOK,
			response:   `{ "success": true }`,
			returnVal:  domain.User{Password: "abc"},
			returnErr:  nil,
		},
		{
			name:       "invalid token",
			request:    verifyPasswordReset,
			statusCode: http.StatusBadRequest,
			response:   utils.ErrorJSON("invalid token"),
			returnVal:  domain.User{},
			returnErr:  domain.ErrInvalidToken,
		},
		{
			name:       "internal error",
			request:    verifyPasswordReset,
			statusCode: http.StatusInternalServerError,
			response:   utils.ErrorJSON("internal server error"),
			returnVal:  domain.User{},
			returnErr:  errors.New("random"),
		},
		{
			name:       "token is missing",
			request:    `{ "token": "", "new_password": "testee" }`,
			statusCode: http.StatusBadRequest,
			response:   utils.ErrorJSON("invalid token"),
			skipMock:   true,
		},
		{
			name:       "token is missing",
			request:    `{ "token": "abc", "new_password": "tst" }`,
			statusCode: http.StatusBadRequest,
			response:   utils.ErrorJSON("incorrect password"),
			skipMock:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			req := utils.RequestBuilder("POST", "/passwordreset")(tc.request)
			rr := httptest.NewRecorder()

			baseExpector := mocks.NewDomain(t)
			domain := baseExpector.EXPECT()

			tokenMatcher := mock.MatchedBy(func(token string) bool {
				return token == "abc"
			})

			passwordMatcher := mock.MatchedBy(func(password string) bool {
				return password == "testee"
			})

			if !tc.skipMock {
				domain.VerifyPasswordReset(mock.Anything, tokenMatcher, passwordMatcher).
					Return(tc.returnVal, tc.returnErr)
			}

			api := NewApi(utils.Config{}, http.NewServeMux(), baseExpector, logger)

			api.postVerifyPasswordReset(rr, req)

			require.Equal(t, rr.Code, tc.statusCode)
			require.JSONEq(t, rr.Body.String(), tc.response)
		})
	}
}
