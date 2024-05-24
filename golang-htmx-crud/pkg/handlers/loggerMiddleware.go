package handlers

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewLoggerMiddleware(log *slog.Logger) echo.MiddlewareFunc {
	logValesFunc := func(_ echo.Context, loggerValues middleware.RequestLoggerValues) error {
		if loggerValues.Error == nil {
			log.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
				slog.String("uri", loggerValues.URI),
				slog.String("method", loggerValues.Method),
				slog.Int("status", loggerValues.Status),
			)
		} else {
			log.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
				slog.String("uri", loggerValues.URI),
				slog.String("method", loggerValues.Method),
				slog.Int("status", loggerValues.Status),
				slog.String("err", loggerValues.Error.Error()),
			)
		}

		return nil
	}

	loggerConfig := middleware.RequestLoggerConfig{
		LogStatus:     true,
		LogMethod:     true,
		LogURI:        true,
		LogError:      true,
		HandleError:   true,
		LogValuesFunc: logValesFunc,
	}

	configuredLogger := middleware.RequestLoggerWithConfig(loggerConfig)

	return configuredLogger
}
