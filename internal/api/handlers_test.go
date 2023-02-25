package api

import (
	"context"
	"errors"
	"github.com/djordjev/auth/internal/domain/mocks"
	"github.com/djordjev/auth/internal/domain/types"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostSignup(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name                  string
		payload               io.Reader
		statusCode            int
		mockSignUpReturnUser  *types.User
		mockSignUpReturnError error
		mockSignUpArgCtx      *context.Context
		mockSignUpArgUser     *types.User
	}{
		{
			name:       "invalid request",
			payload:    strings.NewReader("not json"),
			statusCode: http.StatusBadRequest,
		},
		{
			name:                  "internal error on sign up",
			payload:               strings.NewReader(signUpRequest),
			statusCode:            http.StatusInternalServerError,
			mockSignUpArgCtx:      &ctx,
			mockSignUpArgUser:     userFromRequest(0),
			mockSignUpReturnError: errors.New("error"),
			mockSignUpReturnUser:  &types.User{},
		},
		{
			name:                 "success",
			payload:              strings.NewReader(signUpRequest),
			statusCode:           http.StatusOK,
			mockSignUpReturnUser: userFromRequest(31),
			mockSignUpArgCtx:     &ctx,
			mockSignUpArgUser:    userFromRequest(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, "POST", "/signup", test.payload)
			if err != nil {
				require.FailNow(t, "failed to create test request")
			}

			rr := httptest.NewRecorder()
			mux := http.NewServeMux()

			domain := mocks.NewDomain(t)

			if test.mockSignUpArgUser != nil {
				domain.EXPECT().SignUp(*test.mockSignUpArgCtx, *test.mockSignUpArgUser).
					Return(*test.mockSignUpReturnUser, test.mockSignUpReturnError)
			}

			api := NewApi(utils.Config{}, mux, domain)

			api.postSignup(rr, req)

			require.Equal(t, rr.Code, test.statusCode)
		})
	}
}
