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
	ID        pgtype.Uint32      `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UserID    pgtype.Uint32      `db:"userid"`
	Token     pgtype.Text        `db:"token"`
}

type repositoryVerifyAccount struct {
	ctx context.Context
	db  query
}

func (v *repositoryVerifyAccount) Create(token string, userId uint) (verification domain.VerifyAccount, err error) {
	row := v.db.QueryRow(
		v.ctx,
		"insert into verify_accounts (created_at, user_id, token) values ($1, $2, $3)",
		time.Now(), userId, token,
	)

	var id pgtype.Uint32
	err = row.Scan(&id)

	if err != nil {
		err = fmt.Errorf("model VerifyAccount -> unable to create verification token")
		return
	}

	verification.Token = token
	verification.ID = uint(id.Uint32)
	verification.UserID = userId

	return
}

func (v *repositoryVerifyAccount) Verify(token string) (verification domain.VerifyAccount, err error) {
	verifyRequest := VerifyAccount{}

	row := v.db.QueryRow(
		v.ctx,
		"select * from verify_accounts where token = $1",
		token,
	)

	err = row.Scan(
		&verifyRequest.ID,
		&verifyRequest.CreatedAt,
		&verifyRequest.UserID,
		&verifyRequest.Token,
	)

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
		err = fmt.Errorf("failed to delete verification request for user %d %w", verifyRequest.UserID, err)
		return
	}

	verification.ID = uint(verifyRequest.ID.Uint32)
	verification.Token = token
	verification.UserID = uint(verifyRequest.UserID.Uint32)

	return
}

func newRepositoryVerifyAccount(ctx context.Context, db query) *repositoryVerifyAccount {
	return &repositoryVerifyAccount{ctx: ctx, db: db}
}
