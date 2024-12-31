package generic

import (
	"context"

	"github.com/wolftotem4/golava-core/auth"
)

type NullGuard struct{}

func (ng *NullGuard) User() auth.Authenticatable {
	return nil
}

func (ng *NullGuard) SetUser(user auth.Authenticatable) error {
	return nil
}

func (ng *NullGuard) ID() any {
	return nil
}

func (ng *NullGuard) Validate(ctx context.Context, credentials map[string]any) (bool, error) {
	return false, nil
}

func (ng *NullGuard) Check() bool {
	return false
}

func (ng *NullGuard) Guest() bool {
	return true
}

func (ng *NullGuard) HasUser() bool {
	return false
}
