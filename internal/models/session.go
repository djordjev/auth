package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/djordjev/auth/internal/domain"
	modelErrors "github.com/djordjev/auth/internal/models/errors"
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

	if res := s.redis.Expire(s.ctx, key.String(), time.Duration(10*24*time.Hour)); res.Err() != nil {
		err = fmt.Errorf("unable to set expiration to session %s", key)
		return
	}

	session.User = user
	session.ID = key.String()

	return
}

func (s *repositorySession) Get(key string) (user domain.User, err error) {
	cmd := s.redis.HGetAll(s.ctx, key)
	if cmd.Err() != nil {
		err = fmt.Errorf("failed to get key for session %s %w", key, cmd.Err())
		return
	}

	result, err := cmd.Result()
	if err != nil {
		err = fmt.Errorf("unable to get key for session %s %w", key, cmd.Err())
		return
	}

	if len(result) == 0 {
		err = modelErrors.ErrNotFound
		return
	}

	id, err := strconv.Atoi(result["id"])
	if err != nil {
		err = fmt.Errorf("invalid value in session as user id %s %w", result["id"], err)
		return
	}

	verified, err := strconv.ParseBool(result["verified"])
	if err != nil {
		err = fmt.Errorf("invalid value in session as user id %s verified %s %w", result["id"], result["verified"], err)
		return
	}

	user.ID = uint64(id)
	user.Email = result["email"]
	user.Username = result["username"]
	user.Password = result["password"]
	user.Role = result["role"]
	user.Verified = verified

	return
}

func (s *repositorySession) Delete(key string) error {
	cmd := s.redis.HDel(s.ctx, key)

	if cmd.Err() != nil {
		return fmt.Errorf("unable to delete session %s %w", key, cmd.Err())
	}

	return nil
}

func newRepositorySession(ctx context.Context, redis *redis.Client) *repositorySession {
	return &repositorySession{ctx: ctx, redis: redis}
}
