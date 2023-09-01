package models

import (
	"context"
	"fmt"

	"github.com/djordjev/auth/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type repositorySession struct {
	ctx   context.Context
	redis *redis.Client
}

func (s *repositorySession) Create(user domain.User) (session domain.Session, err error) {
	key, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("unable to generate key for session %w", err)
		return
	}

	values := []string{
		"id", fmt.Sprintf("%d", user.ID),
		"email", user.Email,
		"username", user.Username,
		"role", user.Role,
		"verified", fmt.Sprintf("%t", user.Verified),
	}

	cmd := s.redis.HSet(s.ctx, key.String(), values)

	if cmd.Err() != nil {
		err = fmt.Errorf("unable to store session for user %d in redis %w", user.ID, cmd.Err())
		return
	}

	session.User = user
	session.ID = key.String()

	return
}

func newRepositorySession(ctx context.Context, redis *redis.Client) *repositorySession {
	return &repositorySession{ctx: ctx, redis: redis}
}
