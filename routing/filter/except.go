package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/util"
)

// prevent the middleware from being executed on the specified paths
func Except(middleware gin.HandlerFunc, excepts ...string) gin.HandlerFunc {
	for i := range excepts {
		excepts[i] = util.CleanPath(excepts[i])
	}

	return ExceptMatch(middleware, func(path string) bool {
		for _, except := range excepts {
			if path == except {
				return true
			}
		}
		return false
	})
}

func ExceptMatch(middleware gin.HandlerFunc, match func(path string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		i := instance.MustGetInstance(c)
		router := i.App.Base().Router
		path, isRelative := router.RelativePath(c.Request.URL.Path)

		if isRelative {
			if match(path) {
				c.Next()
				return
			}
		}

		middleware(c)
	}
}
