package logger

import (
	"context"
	"log/slog"
	"strings"

	"github.com/edinstance/distributed-aviation-system/services/flights/internal/config"
)

type LogLevelFilterHandler struct {
	handler slog.Handler
	level   slog.Level
}

func (h *LogLevelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *LogLevelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}
	return h.handler.Handle(ctx, record)
}

func (h *LogLevelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogLevelFilterHandler{
		handler: h.handler.WithAttrs(attrs),
		level:   h.level,
	}
}

func (h *LogLevelFilterHandler) WithGroup(name string) slog.Handler {
	return &LogLevelFilterHandler{
		handler: h.handler.WithGroup(name),
		level:   h.level,
	}
}

func getLogLevel(environment string) slog.Level {
	levelStr := config.App.LogLevel
	if levelStr != "" {
		switch strings.ToUpper(levelStr) {
		case "DEBUG":
			return slog.LevelDebug
		case "INFO":
			return slog.LevelInfo
		case "WARN", "WARNING":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		default:
			return slog.LevelInfo
		}
	}

	switch strings.ToLower(environment) {
	case "dev", "development", "local":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
