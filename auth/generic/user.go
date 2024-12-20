package generic

type User struct {
	ID            int     `db:"id" gorm:"id" json:"id"`
	Username      string  `db:"username" gorm:"username" json:"username"`
	Password      string  `db:"password" gorm:"password" json:"-"`
	RememberToken *string `db:"remember_token" gorm:"remember_token" json:"-"`
	CreatedAt     string  `db:"created_at" gorm:"created_at" json:"created_at"`
	UpdatedAt     string  `db:"updated_at" gorm:"updated_at" json:"updated_at"`
}

func (u *User) GetAuthIdentifierName() string {
	return "id"
}

func (u *User) GetAuthIdentifier() interface{} {
	return u.ID
}

func (u *User) GetAuthPasswordName() string {
	return "password"
}

func (u *User) GetAuthPassword() string {
	return u.Password
}

func (u *User) GetRememberToken() string {
	if u.RememberToken == nil {
		return ""
	}

	return *u.RememberToken
}

func (u *User) SetRememberToken(token string) {
	u.RememberToken = &token
}

func (u *User) GetRememberTokenName() string {
	return "remember_token"
}
