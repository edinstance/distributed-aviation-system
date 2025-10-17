package context

import "github.com/google/uuid"

type UserContext struct {
	UserID uuid.UUID
	OrgID  uuid.UUID
	Roles  string
}
