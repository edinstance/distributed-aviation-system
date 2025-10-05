package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	logexport "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var Logger *slog.Logger

func Init(environment string) (*log.LoggerProvider, error) {
	exporter, err := logexport.New(context.Background(),
		logexport.WithEndpoint(config.App.OtlpGrpcUrl),
		logexport.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("flights-service"),
			semconv.DeploymentEnvironmentName(environment),
		)),
	)

	otelHandler := otelslog.NewHandler("flights-service", otelslog.WithLoggerProvider(provider))
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(environment),
	})

	multiHandler := &MultiLogHandler{
		handlers: []slog.Handler{otelHandler, stdoutHandler},
	}

	levelHandler := &LogLevelFilterHandler{
		handler: multiHandler,
		level:   getLogLevel(environment),
	}

	Logger = slog.New(levelHandler)
	slog.SetDefault(Logger)

	return provider, nil
}

// init initialises the package-level Logger with a JSON stdout handler at Info
// level if Logger has not already been configured. It also registers the created
// logger as the default slog logger.
func init() {
	if Logger == nil {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevel(config.App.Environment),
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

// DebugContext logs a message at debug level using the package Logger.
// Optional args are passed through as key/value attributes to the underlying slog.Logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Logger.DebugContext(ctx, msg, args...)
}

// ErrorContext logs an error-level record associated with the provided context.
// It forwards to the package-level Logger's ErrorContext.
// The variadic args are optional key/value pairs as accepted by slog (typically alternating key, value).
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}
