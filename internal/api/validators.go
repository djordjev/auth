package api

import "fmt"

const minPasswordLength = 5

func validateSignup(request SignUpRequest) error {
	if request.Email == "" {
		return fmt.Errorf("missing email")
	}

	if len(request.Password) < minPasswordLength {
		return fmt.Errorf("password must have at least %d characters", minPasswordLength)
	}

	return nil
}

func validateLogin(request LogInRequest) error {
	if request.Password == "" {
		return fmt.Errorf("missing password")
	}

	if len(request.Password) < minPasswordLength {
		return fmt.Errorf("incorrect password")
	}

	if request.Username == "" && request.Email == "" {
		return fmt.Errorf("missing email or username")
	}

	return nil
}

func validateDeleteAccount(request DeleteAccountRequest) error {
	if request.Username == "" && request.Email == "" {
		return fmt.Errorf("missing email or username")
	}

	if request.Password == "" {
		return fmt.Errorf("missing password")
	}

	return nil
}

func validateForgetPassword(request ForgetPasswordRequest) error {
	if request.Email == "" && request.Username == "" {
		return fmt.Errorf("missing username and password")
	}

	return nil
}

func validateVerifyPasswordResetRequest(request VerifyPasswordResetRequest) error {
	if request.Token == "" {
		return fmt.Errorf("invalid token")
	}

	if len(request.NewPassword) < minPasswordLength {
		return fmt.Errorf("incorrect password")
	}

	return nil
}
