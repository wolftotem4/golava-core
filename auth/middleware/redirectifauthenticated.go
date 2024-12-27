package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
)

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
