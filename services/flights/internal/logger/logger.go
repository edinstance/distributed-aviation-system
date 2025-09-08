package logger

import (
	"context"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(environment string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if environment == "dev" {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

func init() {
	if Logger == nil {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		Logger = slog.New(handler)
		slog.SetDefault(Logger)
	}
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
