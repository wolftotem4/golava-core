package middleware

import (
	"context"
	"crypto/subtle"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/auth/generic"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/session"
)

const sessionHashName = "password_hash"

func AuthenticateSession(c *gin.Context) {
	i := instance.MustGetInstance(c)

	err := authenticateSession(c, i)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Next()

	if i.Auth.Check() {
		storePasswordHashInSession(i.Session.Store, sessionHashName, i.Auth.(auth.StatefulGuard).UserHash())
	}
}

func authenticateSession(ctx context.Context, i *instance.Instance) error {
	if i.Session == nil {
		return nil
	}

	guard, ok := i.Auth.(*generic.SessionGuard)
	if !ok {
		return errors.New("auth guard is not a session guard")
	}

	if !guard.Check() || guard.User().GetAuthPassword() == "" {
		return nil
	}

	authPassword := []byte(guard.User().GetAuthPassword())

	if guard.ViaRemember() {
		recaller, err := guard.GetRecaller()
		if err != nil {
			return err
		}
		passwordHash := []byte(recaller.Hash())

		if len(passwordHash) == 0 || subtle.ConstantTimeCompare(passwordHash, authPassword) != 1 {
			guard.LogoutCurrentDevice(ctx)
			i.Session.Store.Forget(sessionHashName)
			return nil
		}
	}

	sessionHash, exists := getPasswordHashFromSession(i.Session.Store, sessionHashName)
	if exists && subtle.ConstantTimeCompare(sessionHash, authPassword) != 1 {
		guard.LogoutCurrentDevice(ctx)
		i.Session.Store.Forget(sessionHashName)
	}

	return nil
}

func storePasswordHashInSession(store *session.Store, name string, password string) {
	store.Put(name, password)
}

func getPasswordHashFromSession(store *session.Store, name string) ([]byte, bool) {
	value, ok := store.Get(name)
	if !ok {
		return nil, false
	}
	str, ok := value.(string)
	return []byte(str), ok
}
