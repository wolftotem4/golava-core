package middleware

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/golava"
)

func LoadFuncMap(engine *gin.Engine, app golava.GolavaApp) {
	engine.SetFuncMap(template.FuncMap{
		"url": app.Base().Router.URL,
	})
}
