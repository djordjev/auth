package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func ErrorJSON(errorMessage string) string {
	return fmt.Sprintf(`{"error": "%s"}`, errorMessage)
}

func RequestBuilder(method string, path string) func(request string) *http.Request {
	return func(request string) *http.Request {
		payload := strings.NewReader(request)

		req, err := http.NewRequestWithContext(context.TODO(), method, path, payload)
		if err != nil {
			panic(errors.New("failed to create mock request"))
		}

		return InjectLoggerIntoContext(req, NewSilentLogger())
	}
}
