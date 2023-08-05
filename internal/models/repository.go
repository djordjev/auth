package models

import (
	"context"
	"github.com/djordjev/auth/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Atomic(fn domain.AtomicFn) (err error) {
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

func (r *repository) User(ctx context.Context) domain.RepositoryUser {
	return newRepositoryUser(ctx, r.db)
}

func (r *repository) VerifyAccount(ctx context.Context) domain.RepositoryVerifyAccount {
	return newRepositoryVerifyAccount(ctx, r.db)
}

func (r *repository) ForgetPassword(ctx context.Context) domain.RepositoryForgetPassword {
	return newRepositoryForgetPassword(ctx, r.db)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}
