package session

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/wolftotem4/golava-core/cookie"
)

type CookieSessionHandler struct {
	Cookie     cookie.IEncryptableCookieManager
	Expiration time.Duration
}

func (c *CookieSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	value, err := c.Cookie.Encryption().Get(sessionId)
	if errors.Is(err, http.ErrNoCookie) {
		return nil, nil
	}
	return []byte(value), err
}

func (c *CookieSessionHandler) Write(ctx context.Context, sessionId string, data SessionData) error {
	c.Cookie.Encryption().Set(sessionId, string(data.Payload), cookie.WithMaxAge(int(c.Expiration.Seconds())))
	return nil
}

func (c *CookieSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	return 0, nil
}

func (c *CookieSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	c.Cookie.Forget(sessionId)
	return nil
}
