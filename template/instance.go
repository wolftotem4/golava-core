package template

import (
	"fmt"
	"html"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/session"
)

type TemplateDataFunc func(c *gin.Context, data H) H

type H map[string]any

func New(c *gin.Context, funcs ...TemplateDataFunc) H {
	data := make(H)
	for _, f := range funcs {
		data = f(c, data)
	}
	return data
}

func Default(c *gin.Context, funcs ...TemplateDataFunc) H {
	return New(c, append(funcs, WithMetadata, WithErrors, WithOld, WithCsrf, WithAuth)...)
}

func (h H) Wrap(data H) H {
	for k, v := range data {
		h[k] = v
	}
	return h
}

func WithMetadata(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	data["app"] = map[string]interface{}{
		"name": instance.App.Base().Name,
	}
	return data
}

func WithErrors(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	data["errors"] = session.GetFlashErrors(instance.Session.Store)
	return data
}

func WithOld(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	old, ok := instance.Session.Store.GetOldInput()
	if ok {
		data["old"] = old
	} else {
		data["old"] = make(map[string]interface{})
	}
	return data
}

func WithCsrf(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	data["csrf_token"] = instance.Session.Store.Token()
	data["csrf"] = template.HTML(fmt.Sprintf(
		`<input type="hidden" name="_token" value="%s">`,
		html.EscapeString(instance.Session.Store.Token()),
	))

	return data
}

func WithAuth(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	data["auth"] = instance.Auth
	return data
}
