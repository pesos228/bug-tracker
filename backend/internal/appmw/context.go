package appmw

import "context"

type contextKey string

const (
	KeyUserId     = contextKey("userId")
	KeyUserEmail  = contextKey("userEmail")
	KeyUserRoles  = contextKey("userRoles")
	KeyGivenName  = contextKey("givenName")
	KeyFamilyName = contextKey("familyName")
)

func UserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(KeyUserId).(string)
	return userId, ok
}

func UserRolesFromContext(ctx context.Context) ([]string, bool) {
	userRoles, ok := ctx.Value(KeyUserRoles).([]string)
	return userRoles, ok
}

func UserFirstNameFromContext(ctx context.Context) (string, bool) {
	firstName, ok := ctx.Value(KeyGivenName).(string)
	return firstName, ok
}

func UserLastNameFromContext(ctx context.Context) (string, bool) {
	lastName, ok := ctx.Value(KeyFamilyName).(string)
	return lastName, ok
}
