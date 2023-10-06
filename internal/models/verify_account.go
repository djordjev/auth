package models

import (
	"context"
	"fmt"
	"time"

	"github.com/djordjev/auth/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyAccount struct {
	ID        pgtype.Int8        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UserID    pgtype.Int8        `db:"user_id"`
	Token     pgtype.Text        `db:"token"`
}

type repositoryVerifyAccount struct {
	ctx context.Context
	db  query
}

func (v *repositoryVerifyAccount) Create(token string, userId uint64) (verification domain.VerifyAccount, err error) {
	row := v.db.QueryRow(
		v.ctx,
		"insert into verify_accounts (created_at, user_id, token) values ($1, $2, $3) returning id",
		time.Now(), userId, token,
	)

	var id pgtype.Int8
	err = row.Scan(&id)

	if err != nil {
		err = fmt.Errorf("model VerifyAccount -> unable to create verification token")
		return
	}

	verification.Token = token
	verification.ID = uint64(id.Int64)
	verification.UserID = userId

	return
}

func (v *repositoryVerifyAccount) Verify(token string) (verification domain.VerifyAccount, err error) {
	rows, err := v.db.Query(
		v.ctx,
		"select * from verify_accounts where token = $1",
		token,
	)
	if err != nil {
		err = fmt.Errorf("failed to query verify accounts %w", err)
		return
	}

	verifyRequest, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[VerifyAccount])

	if err == pgx.ErrNoRows {
		err = fmt.Errorf("there's no verify request associated with token %s %w", token, err)
		return
	} else if err != nil {
		err = fmt.Errorf("failed to find verify token %s %w", token, err)
		return
	}

	_, err = v.db.Exec(
		v.ctx,
		"delete from verify_accounts where token = $1",
		token,
	)

	if err != nil {
		err = fmt.Errorf("failed to delete verification request for user %d %w", verifyRequest.UserID.Int64, err)
		return
	}

	verification.ID = uint64(verifyRequest.ID.Int64)
	verification.Token = token
	verification.UserID = uint64(verifyRequest.UserID.Int64)

	return
}

func newRepositoryVerifyAccount(ctx context.Context, db query) *repositoryVerifyAccount {
	return &repositoryVerifyAccount{ctx: ctx, db: db}
}
