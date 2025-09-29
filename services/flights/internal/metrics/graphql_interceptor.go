package metrics

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func GraphQLMetricsInterceptor(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)

	opName := rc.OperationName
	if opName == "" {
		opName = "unknown"
	}

	opType := string(rc.Operation.Operation)

	start := time.Now()
	resp := next(ctx)

	return func(ctx context.Context) *graphql.Response {
		res := resp(ctx)
		elapsed := time.Since(start).Seconds()

		success := "success"
		if len(res.Errors) > 0 {
			success = "failure"
		}

		GraphQLRequests.WithLabelValues(opName, opType, success).Inc()
		GraphQLDuration.WithLabelValues(opName, opType).Observe(elapsed)

		return res
	}
}
