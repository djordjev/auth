package models

import (
	"context"

	"github.com/djordjev/auth/internal/domain"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func (r *repository) Atomic(fn domain.AtomicFn) (err error) {
	tx := r.db.Begin()

	newRepo := NewRepository(tx, r.redis)

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

func (r *repository) Session(ctx context.Context) domain.RepositorySession {
	return newRepositorySession(ctx, r.redis)
}

func NewRepository(db *gorm.DB, redis *redis.Client) *repository {
	return &repository{db: db, redis: redis}
}
