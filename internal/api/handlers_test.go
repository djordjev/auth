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
	t.Parallel()

	tests := []struct {
		name                  string
		payload               io.Reader
		statusCode            int
		mockSignUpReturnUser  *types.User
		mockSignUpReturnError error
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
			mockSignUpArgUser:     userFromRequest(0),
			mockSignUpReturnError: errors.New("error"),
			mockSignUpReturnUser:  &types.User{},
		},
		{
			name:                 "success",
			payload:              strings.NewReader(signUpRequest),
			statusCode:           http.StatusOK,
			mockSignUpReturnUser: userFromRequest(31),
			mockSignUpArgUser:    userFromRequest(0),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(context.TODO(), "POST", "/signup", test.payload)
			if err != nil {
				require.FailNow(t, "failed to create test request")
			}

			rr := httptest.NewRecorder()
			mux := http.NewServeMux()

			domain := mocks.NewDomain(t)

			testLogger := utils.NewSilentLogger()
			api := NewApi(utils.Config{}, mux, domain, testLogger)
			reqWithLogger := utils.InjectLoggerIntoContext(req, testLogger)

			if test.mockSignUpArgUser != nil {
				domain.EXPECT().SignUp(reqWithLogger.Context(), *test.mockSignUpArgUser).
					Return(*test.mockSignUpReturnUser, test.mockSignUpReturnError)
			}

			api.postSignup(rr, reqWithLogger)

			require.Equal(t, rr.Code, test.statusCode)
		})
	}
}
