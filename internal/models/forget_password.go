package models

import (
	"context"
	"fmt"
	"time"

	"github.com/djordjev/auth/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ForgetPassword struct {
	ID        pgtype.Int8        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UserID    pgtype.Int8        `db:"user_id"`
	Token     pgtype.Text        `db:"token"`
}

type repositoryForgetPassword struct {
	ctx context.Context
	db  query
}

func (fp *repositoryForgetPassword) Create(token string, userId uint64) (request domain.ForgetPassword, err error) {
	_, err = fp.db.Exec(fp.ctx, "delete from forget_passwords where user_id = $1", userId)

	if err != nil {
		err = fmt.Errorf("failed to delete previous password change requests for user %d %w", userId, err)
		return
	}

	_, err = fp.db.Exec(
		fp.ctx,
		"insert into forget_passwords (created_at, token, user_id) values ($1, $2, $3)",
		time.Now(), token, userId,
	)

	if err != nil {
		err = fmt.Errorf("failed to create password request for user %d %w", userId, err)
		return
	}

	request.Token = token
	request.UserID = userId

	return
}

func (fp *repositoryForgetPassword) Delete(token string) (request domain.ForgetPassword, err error) {

	rows, err := fp.db.Query(
		fp.ctx,
		"select * from forget_passwords where token = $1",
		token,
	)

	if err != nil {
		err = fmt.Errorf("failed to query forget passwords %w", err)
		return
	}

	forgetPasswordReq, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ForgetPassword])

	if err == pgx.ErrNoRows {
		err = fmt.Errorf("there's no reset request associated with token %s %w", token, err)
		return
	} else if err != nil {
		err = fmt.Errorf("failed to find token %s %w", token, err)
		return
	}

	_, err = fp.db.Exec(
		fp.ctx,
		"delete from forget_passwords where token = $1",
		token,
	)

	if err != nil {
		err = fmt.Errorf("failed to delete password reset request for user %d %w", forgetPasswordReq.UserID.Int64, err)
		return
	}

	request.ID = uint64(forgetPasswordReq.ID.Int64)
	request.Token = token
	request.UserID = uint64(forgetPasswordReq.UserID.Int64)

	return
}

func newRepositoryForgetPassword(ctx context.Context, db query) *repositoryForgetPassword {
	return &repositoryForgetPassword{ctx: ctx, db: db}
}
