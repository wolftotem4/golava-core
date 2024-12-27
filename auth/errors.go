package auth

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthenticated = errors.New("unauthenticated")
var ErrPasswordMismatch = errors.New("the given password does not match the current password")
