package foundation

import (
	"context"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/session"
)

func StartSession(c *gin.Context) {
	instance := instance.MustGetInstance(c)
	instance.Session = instance.App.Base().SessionFactory.Make()
	instance.Redirector.Session = instance.Session

	var sessionId string

	migrateId, _ := instance.App.Base().Cookie.Encryption().Get(instance.Session.GetMigrateName())
	if migrateId != "" {
		sessionId = migrateId
	} else {
		sessionId, _ = instance.App.Base().Cookie.Encryption().Get(instance.Session.Name)
	}

	if sessionId != "" {
		instance.Session.Store.ID = sessionId
	}

	err := instance.Session.Store.Start(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	err = collectGarbage(c, instance.Session)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	instance.App.Base().Cookie.Encryption().Set(
		instance.Session.Name,
		instance.Session.Store.ID,
		cookie.WithMaxAge(int(instance.Session.Lifetime.Seconds())),
	)

	if migrateId != "" {
		instance.App.Base().Cookie.Forget(instance.Session.GetMigrateName())
	}

	c.Next()
}

func collectGarbage(ctx context.Context, session *session.SessionManager) error {
	hitLottery := rand.Intn(100) == 0
	if hitLottery {
		_, err := session.Store.Handler.GC(ctx, session.Lifetime)
		return err
	}
	return nil
}
