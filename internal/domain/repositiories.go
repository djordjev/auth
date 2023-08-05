package domain

import "context"

type AtomicFn = func(txRepo Repository) error

type Repository interface {
	Atomic(fn AtomicFn) error
	User(ctx context.Context) RepositoryUser
	VerifyAccount(ctx context.Context) RepositoryVerifyAccount
	ForgetPassword(ctx context.Context) RepositoryForgetPassword
}

type RepositoryUser interface {
	Create(user User) (newUser User, err error)
	Delete(id uint) (success bool, err error)
	GetByEmail(email string) (user User, err error)
	GetByUsername(username string) (user User, err error)
	Verify(user User) error
	SetPassword(user User, password string) error
}

type RepositoryVerifyAccount interface {
	Create(token string, userId uint) (verification VerifyAccount, err error)
	Verify(token string) (verification VerifyAccount, err error)
}

type RepositoryForgetPassword interface {
	Create(userId uint) (request ForgetPassword, err error)
	Delete(token string) (request ForgetPassword, err error)
}