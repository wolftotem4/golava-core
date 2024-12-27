package callback

import "github.com/wolftotem4/golava-core/auth"

type Attempting interface {
	Attempting(name string, credentials map[string]any, remember bool) error
}

type Validated interface {
	Validated(name string, user auth.Authenticatable) error
}

type Login interface {
	Login(name string, user auth.Authenticatable, remember bool) error
}

type Authenticated interface {
	Authenticated(name string, user auth.Authenticatable) error
}

type CurrentDeviceLogout interface {
	CurrentDeviceLogout(name string, user auth.Authenticatable) error
}

type OtherDeviceLogout interface {
	OtherDeviceLogout(name string, user auth.Authenticatable) error
}

type Failed interface {
	Failed(name string, user auth.Authenticatable) error
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

func (e *listener) Attempting(name string, credentials map[string]any, remember bool) error {
	callback, ok := e.base.(Attempting)
	if ok {
		return callback.Attempting(name, credentials, remember)
	}
	return nil
}

func (e *listener) Validated(name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Validated)
	if ok {
		return callback.Validated(name, user)
	}
	return nil
}

func (e *listener) Login(name string, user auth.Authenticatable, remember bool) error {
	callback, ok := e.base.(Login)
	if ok {
		return callback.Login(name, user, remember)
	}
	return nil
}

func (e *listener) Authenticated(name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Authenticated)
	if ok {
		return callback.Authenticated(name, user)
	}
	return nil
}

func (e *listener) CurrentDeviceLogout(name string, user auth.Authenticatable) error {
	callback, ok := e.base.(CurrentDeviceLogout)
	if ok {
		return callback.CurrentDeviceLogout(name, user)
	}
	return nil
}

func (e *listener) OtherDeviceLogout(name string, user auth.Authenticatable) error {
	callback, ok := e.base.(OtherDeviceLogout)
	if ok {
		return callback.OtherDeviceLogout(name, user)
	}
	return nil
}

func (e *listener) Failed(name string, user auth.Authenticatable) error {
	callback, ok := e.base.(Failed)
	if ok {
		return callback.Failed(name, user)
	}
	return nil
}
