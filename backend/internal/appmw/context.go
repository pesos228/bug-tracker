package appmw

import "context"

type contextKey string

const (
	KeyUserId    = contextKey("userId")
	KeyUserEmail = contextKey("userEmail")
	KeyUserRoles = contextKey("userRoles")
)

func UserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(KeyUserId).(string)
	return userId, ok
}

func UserRolesFromContext(ctx context.Context) ([]string, bool) {
	userRoles, ok := ctx.Value(KeyUserRoles).([]string)
	return userRoles, ok
}
