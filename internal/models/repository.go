package models

import (
	"context"

	"gorm.io/gorm"
)

type AtomicFn = func(txRepo Repository) error

//go:generate mockery --name Repository
type Repository interface {
	Atomic(fn AtomicFn) error
	User(ctx context.Context) RepositoryUser
	VerifyAccount(ctx context.Context) RepositoryVerifyAccount
	ForgetPassword(ctx context.Context) RepositoryForgetPassword
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Atomic(fn AtomicFn) (err error) {
	tx := r.db.Begin()

	newRepo := NewRepository(tx)

	err = fn(newRepo)

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return
}

func (r *repository) User(ctx context.Context) RepositoryUser {
	return newRepositoryUser(ctx, r.db)
}

func (r *repository) VerifyAccount(ctx context.Context) RepositoryVerifyAccount {
	return newRepositoryVerifyAccount(ctx, r.db)
}

func (r *repository) ForgetPassword(ctx context.Context) RepositoryForgetPassword {
	return newRepositoryForgetPassword(ctx, r.db)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}
