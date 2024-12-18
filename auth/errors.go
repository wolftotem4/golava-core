package auth

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthenticated = errors.New("unauthenticated")
