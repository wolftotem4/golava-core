package cookie

import "github.com/gin-gonic/gin"

func CookieMiddleware(manager ICookieManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		manager.SetRequest(ctx.Request)
		manager.SetResponseWriter(ctx.Writer)

		ctx.Next()
	}
}
