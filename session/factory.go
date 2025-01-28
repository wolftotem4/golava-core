package session

import "time"

type SessionFactory struct {
	Name     string
	Lifetime time.Duration
	HttpOnly bool
	Handler  SessionHandler
}

func (sm *SessionFactory) Make(sessionId string) *SessionManager {
	if sessionId == "" {
		sessionId = NewSessionId()
	}

	store := NewStore(sessionId, sm.Handler)

	return &SessionManager{
		Name:     sm.Name,
		Store:    store,
		Lifetime: sm.Lifetime,
		HttpOnly: sm.HttpOnly,
	}
}
