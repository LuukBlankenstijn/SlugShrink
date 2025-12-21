package authconfig

import "context"

type UserPermission string

const (
	User      UserPermission = "user"
	SuperUser UserPermission = "superUser"
	Admin     UserPermission = "admin"
)

type baseStrategy struct{}

type authContextKeyType string

const authContextKey authContextKeyType = "authContext"

type authContext struct {
	UserPermission UserPermission
	UserId         *string
}

func (baseStrategy) setContext(ctx context.Context, permission UserPermission, userId *string) context.Context {
	return context.WithValue(ctx, authContextKey, authContext{
		UserPermission: permission,
		UserId:         userId,
	})
}

func (baseStrategy) GetAuthState(ctx context.Context) (*authContext, bool) {
	switch value := ctx.Value(authContextKey).(type) {
	case authContext:
		return &value, true
	default:
		return nil, false
	}
}
