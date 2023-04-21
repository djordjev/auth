package utils

import (
	"context"
	"golang.org/x/exp/slog"
)

func NewSilentLogger() *slog.Logger {
	return slog.New(silent{})
}

type silent struct{}

func (s silent) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (s silent) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (s silent) WithAttrs(_ []slog.Attr) slog.Handler {
	return s
}

func (s silent) WithGroup(_ string) slog.Handler {
	return s
}
