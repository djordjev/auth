package middleware

import (
	"github.com/djordjev/auth/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
)

type Wrapper func(next http.Handler) http.Handler

func Logger(logger *slog.Logger) Wrapper {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			entry := logger.With("method", r.Method, "path", r.URL.Path)

			defer func() {
				status := ww.Status()
				if status >= http.StatusInternalServerError {
					entry.Error("request failed", "status", status, "size", ww.BytesWritten())
				} else {
					entry.Info("request completed", "status", status, "size", ww.BytesWritten())
				}
			}()

			next.ServeHTTP(ww, utils.InjectLoggerIntoContext(r, entry))
		}

		return http.HandlerFunc(fn)
	}
}
