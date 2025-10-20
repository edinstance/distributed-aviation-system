package middleware

import (
	"context"
	"net/http"

	userContext "github.com/edinstance/distributed-aviation-system/services/flights/internal/context"
	"github.com/google/uuid"
)

type contextKey string

const userCtxKey = contextKey("userContext")

func UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header

		userSub := headers.Get("x-user-sub")
		orgID := headers.Get("x-org-id")
		orgName := headers.Get("x-org-name")
		roles := headers.Get("x-user-roles")

		var parsedUserID, parsedOrgID uuid.UUID
		var err error

		if userSub != "" {
			parsedUserID, err = uuid.Parse(userSub)
			if err != nil {
				parsedUserID = uuid.Nil
			}
		}

		if orgID != "" {
			parsedOrgID, err = uuid.Parse(orgID)
			if err != nil {
				parsedOrgID = uuid.Nil
			}
		}

		ctx := context.WithValue(r.Context(), userCtxKey, &userContext.UserContext{
			UserID:  parsedUserID,
			OrgID:   parsedOrgID,
			OrgName: orgName,
			Roles:   roles,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper to get the UserContext inside resolvers or directives.
func GetRequestUserContext(ctx context.Context) *userContext.UserContext {
	val := ctx.Value(userCtxKey)
	if val == nil {
		return &userContext.UserContext{}
	}
	if userCtx, ok := val.(*userContext.UserContext); ok {
		return userCtx
	}
	return &userContext.UserContext{}
}

func SetUserContextInContext(ctx context.Context, userCtx *userContext.UserContext) context.Context {
	return context.WithValue(ctx, userCtxKey, userCtx)
}
