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

var sl = utils.NewSilentLogger()
var mux = http.NewServeMux()

var signUpRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com",
		"role": "admin"
	}
`

var logInRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com"
	}
`

var deleteRequest = `
	{
		"email": "djvukovic@gmail.com",
		"password": "testee"
	}
`

func TestLogIn(t *testing.T) {
	t.Parallel()

	requestBuilder := utils.RequestBuilder("POST", "/login")

	successUser := domain.User{
		ID:       884,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
		Verified: true,
	}

	userMatcher := mock.MatchedBy(func(usr domain.User) bool {
		return usr.Email == "djvukovic@gmail.com" && usr.Password == "testee"
	})

	type testCase struct {
		name         string
		request      *http.Request
		setupDomain  func(*mocks.Domain, *testCase)
		responseCode int
		responseBody string
	}

	tests := []testCase{
		{
			name:    "success",
			request: requestBuilder(logInRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().LogIn(mock.Anything, userMatcher).Return(successUser, nil)
			},
			responseCode: http.StatusOK,
			responseBody: `{"id": 884, "username": "djvukovic", "email": "djvukovic@gmail.com", "role": "admin", "verified": true }`,
		},
		{
			name:         "validation fail",
			request:      requestBuilder(`{}`),
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("missing password"),
		},
		{
			name:         "bad request",
			request:      requestBuilder(``),
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("bad request"),
		},
		{
			name:    "invalid credentials",
			request: requestBuilder(logInRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().LogIn(mock.Anything, userMatcher).Return(domain.User{}, domain.ErrInvalidCredentials)
			},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("invalid credentials"),
		},
		{
			name:    "random error",
			request: requestBuilder(logInRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().LogIn(mock.Anything, userMatcher).Return(domain.User{}, errors.New("random error"))
			},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("failed login attempt"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			rr := httptest.NewRecorder()
			baseMock := mocks.NewDomain(t)

			// Setup mocks
			if tc.setupDomain != nil {
				tc.setupDomain(baseMock, &tc)
			}

			// Run
			api := NewApi(utils.Config{}, mux, baseMock, sl)
			api.postLogin(rr, tc.request)

			// Assertions
			require.Equal(t, rr.Code, tc.responseCode)
			require.JSONEq(t, rr.Body.String(), tc.responseBody)
		})
	}
}

func TestApiDelete(t *testing.T) {
	t.Parallel()

	requestBuilder := utils.RequestBuilder("DELETE", "/account")

	type testCase struct {
		name         string
		request      *http.Request
		setupDomain  func(*mocks.Domain, *testCase)
		responseCode int
		responseBody string
	}

	userMatcher := mock.MatchedBy(func(usr domain.User) bool {
		return usr.Email == "djvukovic@gmail.com" && usr.Password == "testee"
	})

	tests := []testCase{
		{
			name:    "success",
			request: requestBuilder(deleteRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().Delete(mock.Anything, userMatcher).Return(true, nil)
			},
			responseCode: http.StatusOK,
			responseBody: `{ "success": true }`,
		},
		{
			name:         "validation failed",
			request:      requestBuilder(`{ "email": "djvukovic@gmail.com" }`),
			setupDomain:  func(d *mocks.Domain, tc *testCase) {},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("missing password"),
		},
		{
			name:    "delete non existing user",
			request: requestBuilder(deleteRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().Delete(mock.Anything, userMatcher).Return(false, domain.ErrUserNotExist)
			},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("user does not exists"),
		},
		{
			name:    "delete authentication failed",
			request: requestBuilder(deleteRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().Delete(mock.Anything, userMatcher).Return(false, domain.ErrInvalidCredentials)
			},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("authentication failed"),
		},
		{
			name:    "delete error",
			request: requestBuilder(deleteRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().Delete(mock.Anything, userMatcher).Return(false, errors.New("new error"))
			},
			responseCode: http.StatusInternalServerError,
			responseBody: utils.ErrorJSON("internal server error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			rr := httptest.NewRecorder()
			baseMock := mocks.NewDomain(t)

			// Setup mocks
			tc.setupDomain(baseMock, &tc)

			// Run
			api := NewApi(utils.Config{}, mux, baseMock, sl)
			api.deleteAccount(rr, tc.request)

			// Assertions
			require.Equal(t, rr.Code, tc.responseCode)
			require.JSONEq(t, rr.Body.String(), tc.responseBody)
		})
	}
}

func TestApiSignUp(t *testing.T) {
	t.Parallel()

	requestBuilder := utils.RequestBuilder("POST", "/signup")

	newUser := domain.User{
		ID:       31,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}

	type testCase struct {
		name         string
		request      *http.Request
		setupDomain  func(*mocks.Domain, *testCase)
		responseCode int
		responseBody string
	}

	domainUserMatcher := mock.MatchedBy(func(usr domain.User) bool {
		return usr.Username == "djvukovic" &&
			usr.Password == "testee" &&
			usr.Role == "admin" &&
			usr.Email == "djvukovic@gmail.com"
	})

	tests := []testCase{
		{
			name:    "error user already exists",
			request: requestBuilder(signUpRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().SignUp(mock.Anything, domainUserMatcher).Return(domain.User{}, domain.ErrUserAlreadyExists)
			},
			responseCode: http.StatusBadRequest,
			responseBody: utils.ErrorJSON("user with email djvukovic@gmail.com already exists"),
		},
		{
			name:    "success",
			request: requestBuilder(signUpRequest),
			setupDomain: func(d *mocks.Domain, tc *testCase) {
				d.EXPECT().SignUp(mock.Anything, domainUserMatcher).Return(newUser, nil)
			},
			responseCode: http.StatusOK,
			responseBody: `{"id": 31, "username": "djvukovic", "email": "djvukovic@gmail.com", "role": "admin"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			rr := httptest.NewRecorder()
			baseMock := mocks.NewDomain(t)

			// Setup mocks
			tc.setupDomain(baseMock, &tc)

			// Run
			api := NewApi(utils.Config{}, mux, baseMock, sl)
			api.postSignup(rr, tc.request)

			// Assertions
			require.Equal(t, rr.Code, tc.responseCode)
			require.JSONEq(t, rr.Body.String(), tc.responseBody)
		})
	}
}
