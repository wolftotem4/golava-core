package foundation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/instance"
)

func Authenticate(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	if !instance.Auth.Check() {
		c.Error(auth.ErrUnauthenticated)
		c.Abort()
		return
	}

	c.Next()
}

func RedirectIfAuthenticated(redirectTo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := instance.MustGetInstance(c)

		if instance.Auth.Check() {
			instance.Redirector.Redirect(http.StatusSeeOther, redirectTo)
			c.Abort()
			return
		}

		c.Next()
	}
}
