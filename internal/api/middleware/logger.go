package middleware

import (
	"github.com/djordjev/auth/internal/utils"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Wrapper func(next http.Handler) http.Handler

func Logger(logger *zap.SugaredLogger) Wrapper {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			entry := logger.With("path", r.URL.Path)

			defer func() {
				entry.Infow("middleware logger", "status", ww.Status(), "bytes", ww.BytesWritten(), "duration", time.Since(t1))
			}()

			next.ServeHTTP(ww, utils.InjectLoggerIntoContext(r, entry))
		}

		return http.HandlerFunc(fn)
	}
}
