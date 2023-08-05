package utils

import (
	"context"
	"errors"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

const loggerKey = "__app_logger_key"

var ErrNoLoggerInContext = errors.New("no logger in context")

func MustBuildLogger(config Config) *slog.Logger {
	if config.IsDev() {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
		logger := slog.New(handler)

		return logger.With("name", "dev-auth")
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler)

	return logger.With("name", "auth")
}

func MustGetLogger(r *http.Request) *slog.Logger {
	logger, ok := r.Context().Value(loggerKey).(*slog.Logger)

	if !ok {
		panic(ErrNoLoggerInContext)
	}

	return logger
}

func InjectLoggerIntoContext(r *http.Request, entry *slog.Logger) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), loggerKey, entry))
}

func LogError(logger *slog.Logger, err error) {
	logger.Error(err.Error())
}
