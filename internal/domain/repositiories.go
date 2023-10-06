package domain

import (
	"context"
)

type AtomicFn = func(txRepo Repository) error

type Repository interface {
	Atomic(fn AtomicFn) error
	User(ctx context.Context) RepositoryUser
	VerifyAccount(ctx context.Context) RepositoryVerifyAccount
	ForgetPassword(ctx context.Context) RepositoryForgetPassword
	Session(ctx context.Context) RepositorySession
}

type RepositoryUser interface {
	Create(user User) (newUser User, err error)
	Delete(id uint64) (success bool, err error)
	GetByEmail(email string) (user User, err error)
	GetByUsername(username string) (user User, err error)
	Verify(user User) error
	SetPassword(user User, password string) error
}

type RepositoryVerifyAccount interface {
	Create(token string, userId uint64) (verification VerifyAccount, err error)
	Verify(token string) (verification VerifyAccount, err error)
}

type RepositoryForgetPassword interface {
	Create(token string, userId uint64) (request ForgetPassword, err error)
	Delete(token string) (request ForgetPassword, err error)
}

type RepositorySession interface {
	Create(user User) (session Session, err error)
	Get(key string) (user User, err error)
	Delete(key string) error
}
