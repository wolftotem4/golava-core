package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/instance"
)

func CookieMiddleware(factory *cookie.CookieFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var i = instance.MustGetInstance(c)
		i.Cookie = factory.Make(c.Request, c.Writer)
		c.Next()
	}
}
