package directives

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/edinstance/distributed-aviation-system/services/flights/internal/middleware"
	"github.com/google/uuid"
)

func AuthenticationDirective(ctx context.Context, _ any, next graphql.Resolver) (res any, err error) {
	user := middleware.GetRequestUserContext(ctx)

	if user.UserID == uuid.Nil || user.OrgID == uuid.Nil {
		return nil, errors.New("unauthorized: userId or orgId missing")
	}

	if user.OrgName == "" {
		return nil, errors.New("unauthorized: orgName missing")
	}

	return next(ctx)
}
