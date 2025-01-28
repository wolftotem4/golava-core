package middleware

import (
	"context"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/session"
)

func StartSession(factory *session.SessionFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		i := instance.MustGetInstance(c)
		i.Session = factory.Make()
		i.Redirector.Session = i.Session

		var sessionId string

		migrateId, _ := i.Cookie.Encryption().Get(i.Session.GetMigrateName())
		if migrateId != "" {
			sessionId = migrateId
		} else {
			sessionId, _ = i.Cookie.Encryption().Get(i.Session.Name)
		}

		if sessionId != "" {
			i.Session.Store.ID = sessionId
		}

		err := i.Session.Store.Start(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		err = collectGarbage(c, i.Session)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		i.Cookie.Encryption().Set(
			i.Session.Name,
			i.Session.Store.ID,
			cookie.WithMaxAge(int(i.Session.Lifetime.Seconds())),
			cookie.WithHttpOnly(i.Session.HttpOnly),
		)

		if migrateId != "" {
			i.Cookie.Forget(i.Session.GetMigrateName())
		}

		c.Next()
	}
}

func collectGarbage(ctx context.Context, session *session.SessionManager) error {
	hitLottery := rand.Intn(100) == 0
	if hitLottery {
		_, err := session.Store.Handler.GC(ctx, session.Lifetime)
		return err
	}
	return nil
}
