package models

import (
	"github.com/djordjev/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func modelUserToDomainUser(model User) domain.User {
	var username *string
	if !model.Username.Valid {
		username = nil
	} else {
		username = &model.Username.String
	}

	return domain.User{
		ID:       uint(model.ID.Uint32),
		Email:    model.Email.String,
		Username: *username,
		Password: model.Password.String,
		Role:     model.Role.String,
		Verified: model.Verified.Bool,
	}
}

func domainUserToModelUser(user domain.User) User {
	usr := User{
		Email:    pgtype.Text{String: user.Email},
		Password: pgtype.Text{String: user.Password},
		Username: pgtype.Text{String: user.Username, Valid: user.Username != ""},
		Role:     pgtype.Text{String: user.Role},
		Verified: pgtype.Bool{Bool: user.Verified},
	}

	usr.ID = pgtype.Uint32{Uint32: uint32(user.ID)}

	return usr
}
