package models

import (
	"context"

	"github.com/djordjev/auth/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	db    query
	redis *redis.Client
}

func (r *repository) Atomic(fn domain.AtomicFn) (err error) {
	tx, err := r.db.Begin(context.Background())

	newRepo := NewRepository(tx, r.redis)

	err = fn(newRepo)

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
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

func NewRepository(db query, redis *redis.Client) *repository {
	return &repository{db: db, redis: redis}
}

type query interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}
