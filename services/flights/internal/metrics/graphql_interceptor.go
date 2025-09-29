package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/logger"
	"github.com/vektah/gqlparser/v2/ast"
)

func GraphQLMetricsInterceptor(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)

	opName := deriveOpName(rc)
	opType := "unknown"
	if rc != nil && rc.Operation != nil {
		opType = string(rc.Operation.Operation)
	}

	start := time.Now()
	resp := next(ctx)

	return func(ctx context.Context) *graphql.Response {
		res := resp(ctx)
		elapsed := time.Since(start).Seconds()

		result := "success"
		if res != nil && len(res.Errors) > 0 {
			result = "failure"
		}

		GraphQLRequests.WithLabelValues(opName, opType, result).Inc()
		GraphQLDuration.WithLabelValues(opName, opType).Observe(elapsed)

		logger.DebugContext(ctx, "GraphQL operation handled",
			"operation", opName,
			"type", opType,
			"result", result,
			"duration_s", elapsed,
		)

		return res
	}
}

func deriveOpName(rc *graphql.OperationContext) string {
	if rc == nil {
		return "unknown"
	}
	if rc.OperationName != "" {
		return rc.OperationName
	}

	if rc.Operation != nil {
		for _, sel := range rc.Operation.SelectionSet {
			if field, ok := sel.(*ast.Field); ok && field != nil && field.Name != "" {
				return fmt.Sprintf("%s.%s", rc.Operation.Operation, field.Name)
			}
		}
	}
	return "unknown"
}
