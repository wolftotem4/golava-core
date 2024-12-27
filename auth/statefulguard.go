package auth

import "context"

type StatefulGuard interface {
	Guard

	Attempt(ctx context.Context, credentials map[string]any, remember bool, shouldLogin ...ShouldLogin) (bool, error)

	// Log a user into the application without sessions or cookies.
	Once(ctx context.Context, credentials map[string]any) (bool, error)

	// Log the given user ID into the application without sessions or cookies.
	OnceUsingID(ctx context.Context, id any) (bool, error)

	Login(ctx context.Context, user Authenticatable, remember bool) error
	LoginUsingID(ctx context.Context, id any, remember bool) error
	Logout(ctx context.Context) error
	ViaRemember() bool
	UserHash() string
}

type ShouldLogin func(ctx context.Context, user Authenticatable) (valid bool, err error)
