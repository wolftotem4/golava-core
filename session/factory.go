package session

import "time"

type SessionFactory struct {
	Name     string
	Lifetime time.Duration
	HttpOnly bool
	Handler  SessionHandler
}

func (sm *SessionFactory) Make() *SessionManager {
	store := NewStore(NewSessionId(), sm.Handler)

	return &SessionManager{
		Name:     sm.Name,
		Store:    store,
		Lifetime: sm.Lifetime,
		HttpOnly: sm.HttpOnly,
	}
}
