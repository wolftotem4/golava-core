package session

import (
	"context"
	"time"
)

type SessionHandler interface {
	Read(ctx context.Context, sessionId string) ([]byte, error)
	Write(ctx context.Context, sessionId string, data SessionData) error
	GC(ctx context.Context, lifetime time.Duration) (int64, error)
	Destroy(ctx context.Context, sessionId string) error
}
