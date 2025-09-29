package metrics

import (
	"context"
	"errors"
	"strings"
	"time"

	"connectrpc.com/connect"
)

type GrpcMetricsInterceptor struct{}

func (GrpcMetricsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()

		resp, err := next(ctx, req)
		duration := time.Since(start).Seconds()

		status := "OK"
		var connectErr *connect.Error
		if errors.As(err, &connectErr) {
			status = connectErr.Code().String()
		}

		service, method := "unknown", "unknown"
		parts := strings.SplitN(req.Spec().Procedure, "/", 2)
		if len(parts) == 2 {
			service, method = parts[0], parts[1]
		}

		GrpcRequests.WithLabelValues(service, method, status).Inc()
		GrpcDuration.WithLabelValues(service, method).Observe(duration)

		return resp, err
	}
}

func (GrpcMetricsInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (GrpcMetricsInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
