package utils

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

const LoggerKey = "__app_logger_key"

func MustBuildLogger(config Config) *zap.SugaredLogger {
	if config.IsDev() {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		return logger.Sugar().Named("dev-auth")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger.Sugar().Named("auth")
}

func InjectLoggerIntoContext(r *http.Request, logger *zap.SugaredLogger) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), LoggerKey, logger))
}

func MustGetLoggerFromRequest(r *http.Request) *zap.SugaredLogger {
	ctx := r.Context()

	logger := ctx.Value(LoggerKey).(*zap.SugaredLogger)

	if logger == nil {
		panic("no logger in context")
	}

	return logger
}
