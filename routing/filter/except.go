package filter

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
)

// prevent the middleware from being executed on the specified paths
func Except(middleware gin.HandlerFunc, excepts ...string) gin.HandlerFunc {
	for i := range excepts {
		excepts[i] = cleanPath(excepts[i])
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

// cleanPath normalizes a path by converting backslashes to forward slashes
// and removing leading/trailing slashes
func cleanPath(path string) string {
	if path == "" {
		return ""
	}
	// Convert backslashes to forward slashes for cross-platform compatibility
	cleaned := strings.ReplaceAll(path, "\\", "/")
	// Remove leading and trailing slashes
	return strings.Trim(cleaned, "/")
}
