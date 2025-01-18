package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
)

func CookieMiddleware(c *gin.Context) {
	var (
		i       = instance.MustGetInstance(c)
		a       = i.App.Base()
		manager = a.CookieFactory.Make()
	)

	manager.SetRequest(c.Request)
	manager.SetResponseWriter(c.Writer)

	i.Cookie = manager

	c.Next()
}
