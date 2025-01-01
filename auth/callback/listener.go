package callback

import (
	"context"

	"github.com/wolftotem4/golava-core/auth"
)

type Attempting interface {
	Attempting(ctx context.Context, name string, credentials map[string]any, remember bool) error
}

type Validated interface {
	Validated(ctx context.Context, name string, user auth.Authenticatable) error
}

type Login interface {
	Login(ctx context.Context, name string, user auth.Authenticatable, remember bool) error
}

type Authenticated interface {
	Authenticated(ctx context.Context, name string, user auth.Authenticatable) error
}

type CurrentDeviceLogout interface {
	CurrentDeviceLogout(ctx context.Context, name string, user auth.Authenticatable) error
}

type OtherDeviceLogout interface {
	OtherDeviceLogout(ctx context.Context, name string, user auth.Authenticatable) error
}

type Failed interface {
	Failed(ctx context.Context, name string, user auth.Authenticatable) error
}

type Callbacks interface {
	Attempting
	Validated
	Login
	Authenticated
	CurrentDeviceLogout
	OtherDeviceLogout
	Failed
}

type listener struct {
	base any
}

func Listen(base any) Callbacks {
	callbacks, ok := base.(Callbacks)
	if ok {
		return callbacks
	}

	return &listener{base}
}

func (e *listener) Attempting(ctx context.Context, name string, credentials map[string]any, remember bool) error {
	callback, ok := e.base.(Attempting)
	if ok {
		return callback.Attempting(ctx, name, credentials, remember)
	}
	return nil
}

func (e *listener) Validated(ctx context.Context, name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Validated)
	if ok {
		return callback.Validated(ctx, name, user)
	}
	return nil
}

func (e *listener) Login(ctx context.Context, name string, user auth.Authenticatable, remember bool) error {
	callback, ok := e.base.(Login)
	if ok {
		return callback.Login(ctx, name, user, remember)
	}
	return nil
}

func (e *listener) Authenticated(ctx context.Context, name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Authenticated)
	if ok {
		return callback.Authenticated(ctx, name, user)
	}
	return nil
}

func (e *listener) CurrentDeviceLogout(ctx context.Context, name string, user auth.Authenticatable) error {
	callback, ok := e.base.(CurrentDeviceLogout)
	if ok {
		return callback.CurrentDeviceLogout(ctx, name, user)
	}
	return nil
}

func (e *listener) OtherDeviceLogout(ctx context.Context, name string, user auth.Authenticatable) error {
	callback, ok := e.base.(OtherDeviceLogout)
	if ok {
		return callback.OtherDeviceLogout(ctx, name, user)
	}
	return nil
}

func (e *listener) Failed(ctx context.Context, name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Failed)
	if ok {
		return callback.Failed(ctx, name, user)
	}
	return nil
}
