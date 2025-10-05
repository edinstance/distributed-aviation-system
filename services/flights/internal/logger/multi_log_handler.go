package logger

import (
	"context"
	"log/slog"
)

type MultiLogHandler struct {
	handlers []slog.Handler
}

func (h *MultiLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiLogHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, handler := range h.handlers {
		err := handler.Handle(ctx, record)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (h *MultiLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiLogHandler{handlers: handlers}
}

func (h *MultiLogHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &MultiLogHandler{handlers: handlers}
}
