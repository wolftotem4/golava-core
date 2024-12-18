package auth

import "context"

type Guard interface {
	Check() bool
	Guest() bool
	User() Authenticatable
	ID() any
	Validate(ctx context.Context, credentials map[string]any) (bool, error)
	HasUser() bool
	SetUser(user Authenticatable)
}
