package metrics

import (
	"context"
	"errors"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcMetricsInterceptor struct{}

func (GrpcMetricsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		resp, err := next(ctx, req)
		duration := time.Since(start).Seconds()

		statusCode := "OK"
		var connectErr *connect.Error
		if errors.As(err, &connectErr) {
			statusCode = connectErr.Code().String()
		}

		procedure, service := parseProcedureAndService(req.Spec().Procedure)
		if procedure == "" {
			procedure = "unknown"
		}
		if service == "" {
			service = "unknown"
		}

		GrpcRequests.WithLabelValues("inbound", procedure, service, statusCode).Inc()
		GrpcDuration.WithLabelValues("inbound", procedure, service).Observe(duration)

		logger.DebugContext(ctx, "gRPC request handled",
			"direction", "inbound",
			"procedure", procedure,
			"service", service,
			"status", statusCode,
			"duration_s", duration,
		)

		return resp, err
	}
}

func (GrpcMetricsInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (GrpcMetricsInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func OutboundGrpcUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start).Seconds()

		statusCode := codes.OK.String()
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code().String()
		} else if err != nil {
			statusCode = codes.Unknown.String()
		}

		procedure, service := parseProcedureAndService(method)
		if procedure == "" {
			procedure = "unknown"
		}
		if service == "" {
			service = "unknown"
		}

		GrpcRequests.WithLabelValues("outbound", procedure, service, statusCode).Inc()
		GrpcDuration.WithLabelValues("outbound", procedure, service).Observe(duration)

		logger.DebugContext(ctx, "gRPC request handled",
			"direction", "outbound",
			"procedure", procedure,
			"service", service,
			"status", statusCode,
			"duration_s", duration,
		)

		return err
	}
}

func parseProcedureAndService(raw string) (procedure string, service string) {
	procedure = strings.TrimPrefix(raw, "/")
	service = "unknown"

	if procedure != "" {
		if i := strings.IndexByte(procedure, '/'); i >= 0 {
			service = procedure[:i]
		} else {
			service = procedure
		}
	}
	return
}
