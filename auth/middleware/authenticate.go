package middleware

import (
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
