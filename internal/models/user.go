package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"
	modelErrors "github.com/djordjev/auth/internal/models/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        pgtype.Int8        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	Email     pgtype.Text        `db:"email"`
	Password  pgtype.Text        `db:"password"`
	Username  pgtype.Text        `db:"username"`
	Role      pgtype.Text        `db:"role"`
	Verified  pgtype.Bool        `db:"verified"`
	Payload   []byte             `db:"payload"`
}

type repositoryUser struct {
	ctx context.Context
	db  query
}

func (r *repositoryUser) GetByEmail(email string) (user domain.User, err error) {
	rows, err := r.db.Query(r.ctx, "select * from users where email = $1", email)

	if err != nil {
		err = fmt.Errorf("model GetByEmail -> can not execute query %w", err)
		return
	}

	modelUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])

	if err == pgx.ErrNoRows {
		err = modelErrors.ErrNotFound
		return
	} else if err != nil {
		err = fmt.Errorf("model GetByEmail -> find user by email %s, %w", email, err)
		return
	}

	user = modelUserToDomainUser(modelUser)

	return
}

func (r *repositoryUser) GetByUsername(username string) (user domain.User, err error) {
	rows, err := r.db.Query(r.ctx, "select * from users where username = $1", username)
	if err != nil {
		err = fmt.Errorf("models GetByUsername -> unable to execute query")
		return
	}

	modelUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])

	if err == pgx.ErrNoRows {
		err = modelErrors.ErrNotFound
		return
	} else if err != nil {
		err = fmt.Errorf("model GetByUsername -> find user by username %s, %w", username, err)
		return
	}

	user = modelUserToDomainUser(modelUser)

	return
}

func (r *repositoryUser) Create(user domain.User) (newUser domain.User, err error) {
	modelUser := domainUserToModelUser(user)

	result := r.db.QueryRow(
		r.ctx,
		"insert into users (email, password, username, role, verified, payload) values ($1, $2, $3, $4, $5, $6) returning id",
		modelUser.Email, modelUser.Password, modelUser.Username, modelUser.Role, modelUser.Verified, modelUser.Payload,
	)

	var id pgtype.Int8
	err = result.Scan(&id)

	if err != nil {
		err = fmt.Errorf("model Create -> unable to create user %w", err)
		return
	}

	newUser = modelUserToDomainUser(modelUser)
	newUser.ID = uint64(id.Int64)

	return
}

func (r *repositoryUser) Delete(id uint64) (success bool, err error) {

	result, err := r.db.Exec(r.ctx, "delete from users where id = $1", id)

	if result.RowsAffected() != 1 {
		err = fmt.Errorf("model Delete -> user with id %d does not exist %w", id, err)
		return
	}

	if err != nil {
		err = fmt.Errorf("model Delete -> failed to delete user with id %d %w", id, err)
		return
	}

	return true, nil
}

func (r *repositoryUser) Verify(user domain.User) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	result, err := r.db.Exec(r.ctx, "update users set verified = true where id = $1", user.ID)

	if err != nil {
		return fmt.Errorf("failed to verify user with ID %d %w", user.ID, err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("user with id %d does not exist", user.ID)
	}

	return nil
}

func (r *repositoryUser) SetPassword(user domain.User, password string) error {
	if user.ID == 0 {
		return fmt.Errorf("missing user ID in update function")
	}

	result, err := r.db.Exec(r.ctx, "update users set password = $1 where id = $2", password, user.ID)

	if err != nil {
		return fmt.Errorf("failed to set password for user with ID %d %w", user.ID, err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("user with id %d does not exist", user.ID)
	}

	return nil
}

func newRepositoryUser(ctx context.Context, db query) *repositoryUser {
	return &repositoryUser{ctx: ctx, db: db}
}
