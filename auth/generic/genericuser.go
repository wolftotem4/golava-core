package generic

type GenericUser struct {
	ID            int     `db:"id" gorm:"id" json:"id"`
	Name          string  `db:"name" gorm:"name" json:"name"`
	Password      string  `db:"password" gorm:"password" json:"-"`
	RememberToken *string `db:"remember_token" gorm:"remember_token" json:"-"`
	CreatedAt     string  `db:"created_at" gorm:"created_at" json:"created_at"`
	UpdatedAt     string  `db:"updated_at" gorm:"updated_at" json:"updated_at"`
}

func (u *GenericUser) GetAuthIdentifierName() string {
	return "id"
}

func (u *GenericUser) GetAuthIdentifier() interface{} {
	return u.ID
}

func (u *GenericUser) GetAuthPasswordName() string {
	return "password"
}

func (u *GenericUser) GetAuthPassword() string {
	return u.Password
}

func (u *GenericUser) GetRememberToken() string {
	if u.RememberToken == nil {
		return ""
	}

	return *u.RememberToken
}

func (u *GenericUser) SetRememberToken(token string) {
	u.RememberToken = &token
}

func (u *GenericUser) GetRememberTokenName() string {
	return "remember_token"
}
