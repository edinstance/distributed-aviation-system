package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// TracingHandler wraps an slog.Handler and automatically adds OpenTelemetry
// tracing information (trace_id, span_id, service_name) to log records when available.
type TracingHandler struct {
	handler     slog.Handler
	serviceName string
}

// NewTracingHandler creates a new TracingHandler that wraps the provided handler
// and adds OpenTelemetry tracing context to log records.
func NewTracingHandler(handler slog.Handler, serviceName string) *TracingHandler {
	return &TracingHandler{
		handler:     handler,
		serviceName: serviceName,
	}
}

// Handle processes the log record by adding tracing information and forwarding to the wrapped handler.
func (h *TracingHandler) Handle(ctx context.Context, record slog.Record) error {
	// Extract span context from the provided context
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		// Add trace_id if available
		if spanCtx.HasTraceID() {
			record.AddAttrs(slog.String("trace_id", spanCtx.TraceID().String()))
		}
		// Add span_id if available
		if spanCtx.HasSpanID() {
			record.AddAttrs(slog.String("span_id", spanCtx.SpanID().String()))
		}
	}

	// Always add service_name
	record.AddAttrs(slog.String("service_name", h.serviceName))

	return h.handler.Handle(ctx, record)
}

// Enabled forwards the call to the wrapped handler.
func (h *TracingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs forwards the call to the wrapped handler.
func (h *TracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewTracingHandler(h.handler.WithAttrs(attrs), h.serviceName)
}

// WithGroup forwards the call to the wrapped handler.
func (h *TracingHandler) WithGroup(name string) slog.Handler {
	return NewTracingHandler(h.handler.WithGroup(name), h.serviceName)
}