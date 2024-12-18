package auth

type Authenticatable interface {
	GetAuthIdentifierName() string
	GetAuthIdentifier() any
	GetAuthPasswordName() string
	GetAuthPassword() string
	GetRememberToken() string
	SetRememberToken(string)
	GetRememberTokenName() string
}
