package session

import (
	"time"
)

type SessionManager struct {
	Name     string
	Store    *Store
	Lifetime time.Duration
	HttpOnly bool
}

func (sm *SessionManager) GetMigrateName() string {
	return sm.Name + "_migrate"
}
