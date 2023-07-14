package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSignup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request SignUpRequest
		err     string
	}{
		{
			name:    "success",
			request: SignUpRequest{Email: "djvukovic@gmail.com", Password: "testee"},
		},
		{
			name:    "missing email",
			request: SignUpRequest{Password: "testee"},
			err:     "missing email",
		},
		{
			name:    "incorrect password",
			request: SignUpRequest{Email: "djvukovic@gmail.com", Password: "tst"},
			err:     "password must have at least 5 characters",
		},
		{
			name:    "incorrect email",
			request: SignUpRequest{Email: "incorrectmail", Password: "testee"},
			err:     "incorrect email address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSignup(tc.request)
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request LogInRequest
		err     string
	}{
		{
			name:    "success email",
			request: LogInRequest{Password: "testee", Email: "djvukovic@gmail.com"},
		},
		{
			name:    "success username",
			request: LogInRequest{Password: "testee", Username: "djvukovic"},
		},
		{
			name:    "error password",
			request: LogInRequest{Email: "djvukovic@gmail.com"},
			err:     "missing password",
		},
		{
			name:    "incorrect password",
			request: LogInRequest{Email: "djvukovic@gmail.com", Password: "tst"},
			err:     "incorrect password",
		},
		{
			name:    "missing username and email",
			request: LogInRequest{Password: "testee"},
			err:     "missing email or username",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLogin(tc.request)

			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDeleteAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request DeleteAccountRequest
		err     string
	}{
		{
			name:    "success username",
			request: DeleteAccountRequest{Username: "djvukovic", Password: "testee"},
		},
		{
			name:    "success email",
			request: DeleteAccountRequest{Email: "djvukovic@gmail.com", Password: "testee"},
		},
		{
			name:    "missing username and password",
			request: DeleteAccountRequest{Password: "testee"},
			err:     "missing email or username",
		},
		{
			name:    "missing password",
			request: DeleteAccountRequest{Email: "djvukovic@gmail.com"},
			err:     "missing password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDeleteAccount(tc.request)

			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateForgetPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request ForgetPasswordRequest
		err     string
	}{
		{
			name:    "success email",
			request: ForgetPasswordRequest{Email: "djvukovic@gmail.com"},
		},
		{
			name:    "success username",
			request: ForgetPasswordRequest{Username: "djvukovic"},
		},
		{
			name:    "error missing username and password",
			request: ForgetPasswordRequest{},
			err:     "missing username and password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateForgetPassword(tc.request)

			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateVerifyPasswordReset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request VerifyPasswordResetRequest
		err     string
	}{
		{
			name:    "success",
			request: VerifyPasswordResetRequest{Token: "token", NewPassword: "testee"},
		},
		{
			name:    "missing token",
			request: VerifyPasswordResetRequest{NewPassword: "testee"},
			err:     "invalid token",
		},
		{
			name:    "incorrect password",
			request: VerifyPasswordResetRequest{Token: "token", NewPassword: "tst"},
			err:     "incorrect password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateVerifyPasswordResetRequest(tc.request)

			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
