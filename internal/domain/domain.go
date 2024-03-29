package domain

import (
	"context"
	"errors"
	"fmt"

	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"github.com/djordjev/auth/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Setup struct {
	ctx    context.Context
	logger *slog.Logger
}

func NewSetup(ctx context.Context, logger *slog.Logger) Setup {
	return Setup{
		ctx:    ctx,
		logger: logger,
	}
}

type Domain interface {
	SignUp(setup Setup, user User) (newUser User, err error)
	LogIn(setup Setup, user User) (existing User, sessionKey string, err error)
	Delete(setup Setup, user User) (deleted bool, err error)
	VerifyAccount(setup Setup, token string) (verified bool, err error)
	ResetPasswordRequest(setup Setup, user User) (sentTo User, err error)
	VerifyPasswordReset(setup Setup, token string, password string) (updated User, err error)
	Session(setup Setup, token string) (user User, err error)
	Logout(setup Setup, token string) (err error)
}

func NewDomain(repository Repository, config utils.Config, notifier Notifier) Domain {
	return &domain{db: repository, config: config, notifier: notifier}
}

type domain struct {
	db       Repository
	config   utils.Config
	notifier Notifier
}

func (d *domain) LogIn(setup Setup, user User) (existingUser User, sessionKey string, err error) {
	userModel := d.db.User(setup.ctx)

	if user.Username != "" {
		existingUser, err = userModel.GetByUsername(user.Username)
	} else {
		existingUser, err = userModel.GetByEmail(user.Email)
	}

	if err != nil {
		err = ErrUserNotExist
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		existingUser = User{}
		err = ErrInvalidCredentials
		return
	}

	// Create session
	session, err := d.db.Session(setup.ctx).Create(existingUser)
	if err != nil {
		err = fmt.Errorf("unable to create session for user id %d %w", existingUser.ID, err)
		return
	}

	sessionKey = session.ID

	return

}

func (d *domain) Delete(setup Setup, user User) (deleted bool, err error) {
	userModel := d.db.User(setup.ctx)

	var existingUser User
	if user.Username != "" {
		existingUser, err = userModel.GetByUsername(user.Username)
	} else {
		existingUser, err = userModel.GetByEmail(user.Email)
	}

	if err != nil {
		err = ErrUserNotExist
		deleted = false
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		deleted = false
		err = ErrInvalidCredentials
		return
	}

	deleted, err = userModel.Delete(existingUser.ID)
	return
}

func (d *domain) ResetPasswordRequest(setup Setup, user User) (sentTo User, err error) {
	forgetPasswordModel := d.db.ForgetPassword(setup.ctx)
	userModel := d.db.User(setup.ctx)

	if user.Email != "" {
		sentTo, err = userModel.GetByEmail(user.Email)
	} else {
		sentTo, err = userModel.GetByUsername(user.Username)
	}

	if err == modelErrors.ErrNotFound {
		err = ErrUserNotExist
		return
	} else if err != nil {
		err = fmt.Errorf("domain ResetPasswordRequest -> failed to fetch user")
		return
	}

	token, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("domain ResetPasswordRequest -> failed to generate token %w", err)
		return
	}

	_, err = forgetPasswordModel.Create(token.String(), sentTo.ID)

	if err != nil {
		sentTo = User{}
		return
	}

	return
}

func (d *domain) SignUp(setup Setup, user User) (newUser User, err error) {
	userModel := d.db.User(setup.ctx)

	_, err = userModel.GetByEmail(user.Email)
	if err == nil {
		err = ErrUserAlreadyExists
		return
	}

	if !errors.Is(err, modelErrors.ErrNotFound) {
		err = fmt.Errorf("domain SignUp -> error looking for user %s: %w", user.Email, err)
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		err = fmt.Errorf("domain SignUp -> failed to generate password hash %w", err)
		return
	}

	user.Password = string(hash)
	user.Verified = !d.config.RequireVerification

	err = d.db.Atomic(func(txRepo Repository) error {
		newUserCopy, err := txRepo.User(setup.ctx).Create(user)
		if err != nil {
			err = fmt.Errorf("domain SignUp -> failed to create a new user %w", err)
			return err
		}

		newUser = newUserCopy

		if !d.config.RequireVerification {
			return nil
		}

		token, e := uuid.NewUUID()
		if e != nil {
			return e
		}

		verificationReq, e := txRepo.VerifyAccount(setup.ctx).Create(token.String(), newUserCopy.ID)
		if e != nil {
			return e
		}

		verificationLink := fmt.Sprintf("%s?t=%s", d.config.VerificationLink, verificationReq.Token)
		text := fmt.Sprintf(verifyTextTemplate, user.Email, verificationLink)
		html := fmt.Sprintf(verifyHtmlTemplate, user.Email, verificationLink)

		e = d.notifier.Send(user.Email, "Verify new account", text, html)
		if e != nil {
			return e
		}

		newUser = newUserCopy
		return nil
	})

	return
}

func (d *domain) Session(setup Setup, token string) (user User, err error) {
	user, err = d.db.Session(setup.ctx).Get(token)

	if err == modelErrors.ErrNotFound {
		err = ErrNoSession
		return
	} else if err != nil {
		err = fmt.Errorf("unable to get session %s %w", token, err)
		return
	}

	return
}

func (d *domain) Logout(setup Setup, token string) (err error) {
	err = d.db.Session(setup.ctx).Delete(token)
	if err != nil {
		err = fmt.Errorf("unable to log out user for token %s %w", token, err)
	}

	return
}
