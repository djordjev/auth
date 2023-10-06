package models

import (
	"context"

	"github.com/djordjev/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func storeUser(user domain.User) (domain.User, error) {
	result := dbConnection.QueryRow(
		context.Background(),
		"insert into users (email, password, username, role) values ($1, $2, $3, $4) returning id",
		user.Email, user.Password, user.Username, user.Role,
	)

	var id pgtype.Int8
	if err := result.Scan(&id); err != nil {
		return user, err
	}

	user.ID = uint64(id.Int64)

	return user, nil
}
