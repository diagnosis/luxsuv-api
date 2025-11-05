package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const correlationIDKey ctxKey = "correlation_id"

var globalLogger *slog.Logger

func init() {
	env := os.Getenv("APP_ENV")
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	globalLogger = slog.New(handler)
}

func Get() *slog.Logger {
	return globalLogger
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := globalLogger
	if id := GetCorrelationID(ctx); id != "" {
		logger = logger.With("correlation_id", id)
	}
	return logger
}

func Info(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).InfoContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).ErrorContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).DebugContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).WarnContext(ctx, msg, args...)
}
