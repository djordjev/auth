package models

import (
	"encoding/json"

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

	domainUser := domain.User{
		ID:       uint64(model.ID.Int64),
		Email:    model.Email.String,
		Username: *username,
		Password: model.Password.String,
		Role:     model.Role.String,
		Verified: model.Verified.Bool,
	}

	var payload map[string]any
	if model.Payload != nil {
		payload = make(map[string]any)
		json.Unmarshal(model.Payload, &payload)
		domainUser.Payload = payload
	}

	return domainUser
}

func domainUserToModelUser(user domain.User) User {
	usr := User{
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: pgtype.Text{String: user.Password, Valid: true},
		Username: pgtype.Text{String: user.Username, Valid: user.Username != ""},
		Role:     pgtype.Text{String: user.Role, Valid: true},
		Verified: pgtype.Bool{Bool: user.Verified, Valid: true},
	}

	if user.Payload != nil {
		bytes, _ := json.Marshal(user.Payload)
		usr.Payload = bytes
	}

	usr.ID = pgtype.Int8{Int64: int64(user.ID), Valid: true}

	return usr
}
