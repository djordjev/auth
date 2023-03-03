package models

import (
	"context"
	"gorm.io/gorm"
)

type AtomicFn = func(txRepo Repository) error

//go:generate mockery --name Repository
type Repository interface {
	Atomic(fn AtomicFn)
	User(ctx context.Context) RepositoryUser
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Atomic(fn AtomicFn) {
	tx := r.db.Begin()

	newRepo := NewRepository(tx)

	err := fn(newRepo)

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
}

func (r *repository) User(ctx context.Context) RepositoryUser {
	return newRepositoryUser(ctx, r.db)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}
