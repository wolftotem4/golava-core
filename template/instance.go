package template

import (
	"fmt"
	"html"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/lang"
	"github.com/wolftotem4/golava-core/session"
)

type TemplateDataFunc func(c *gin.Context, data H) H

type H map[string]any

var DefaultFuncs = []TemplateDataFunc{
	WithMetadata, WithErrors, WithOld,
	WithCsrf, WithAuth, WithTranslator,
}

func New(c *gin.Context, funcs ...TemplateDataFunc) H {
	data := make(H)
	for _, f := range funcs {
		data = f(c, data)
	}
	return data
}

func Default(c *gin.Context, funcs ...TemplateDataFunc) H {
	return New(c, append(funcs, DefaultFuncs...)...)
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

	if instance.Session != nil {
		data["errors"] = session.GetFlashErrors(instance.Session.Store)
	}

	return data
}

func WithOld(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	if instance.Session != nil {
		old, ok := instance.Session.Store.GetOldInput()
		if ok {
			data["old"] = old
			return data
		}
	}

	data["old"] = make(map[string]interface{})
	return data
}

func WithCsrf(c *gin.Context, data H) H {
	instance := instance.MustGetInstance(c)

	if instance.Session != nil {
		data["csrf_token"] = instance.Session.Store.Token()
		data["csrf"] = template.HTML(fmt.Sprintf(
			`<input type="hidden" name="_token" value="%s">`,
			html.EscapeString(instance.Session.Store.Token()),
		))
	}

	return data
}

func WithAuth(c *gin.Context, data H) H {
	data["auth"] = instance.MustGetInstance(c).Auth
	return data
}

func WithTranslator(c *gin.Context, data H) H {
	i := instance.MustGetInstance(c)
	fallback := i.App.Base().Translation.GetFallback()
	data["T"] = i.GetUserPreferredTranslator(lang.Fallback(fallback), lang.Soft)
	return data
}

func PassFlash(keys ...string) TemplateDataFunc {
	return func(c *gin.Context, data H) H {
		instance := instance.MustGetInstance(c)

		for _, key := range keys {
			data[key] = instance.Session.Store.Attributes[key]
		}
		return data
	}
}
