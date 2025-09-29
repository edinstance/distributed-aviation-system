package logger

import (
	"context"
	"log/slog"
	"os"
)

var Logger *slog.Logger

// Init configures the package logger based on the given environment.
//
// If environment == "dev" the logger level is set to debug; otherwise it defaults
// to info. The logger writes JSON-formatted records to stdout. This function
// sets the package-level Logger and registers it as the default slog logger.
//
// The environment parameter accepts the string "dev" to enable debug logging.
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

// init initialises the package-level Logger with a JSON stdout handler at Info
// level if Logger has not already been configured. It also registers the created
// logger as the default slog logger.
func init() {
	if Logger == nil {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		Logger = slog.New(handler)
		slog.SetDefault(Logger)
	}
}

// Info logs an informational message using the package-level Logger.
// msg is the message text; args are optional key/value pairs passed to slog.
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// Debug logs a message at debug level using the package Logger.
// Optional args are passed through as key/value attributes to the underlying slog.Logger.
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

// Error logs an error-level message using the package-level Logger.
// It accepts a message and optional key/value pairs (slog-style args) and forwards them to Logger.Error.
func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

// Warn logs msg at the warning level.
// Additional args can be provided as key/value pairs and are forwarded to the package Logger.
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// InfoContext logs an info-level message using the provided context.
//
// The context `ctx` is passed to the underlying logger and can be used to
// include contextual values. `args` are optional key-value pairs (as accepted
// by slog) to attach structured fields to the log entry.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

// ErrorContext logs an error-level record associated with the provided context.
// It forwards to the package-level Logger's ErrorContext.
// The variadic args are optional key/value pairs as accepted by slog (typically alternating key, value).
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
