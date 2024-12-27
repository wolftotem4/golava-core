package auth

import "context"

type UserProvider interface {
	RetrieveById(ctx context.Context, identifier any) (Authenticatable, error)
	RetrieveByToken(ctx context.Context, identifier any, token string) (Authenticatable, error)
	UpdateRememberToken(ctx context.Context, user Authenticatable, token string) error
	RetrieveByCredentials(ctx context.Context, credentials map[string]any) (Authenticatable, error)
	ValidateCredentials(ctx context.Context, user Authenticatable, credentials map[string]any) (bool, error)
	RehashPasswordIfRequired(ctx context.Context, user Authenticatable, credentials map[string]any, force bool) (newhash string, err error)
}
