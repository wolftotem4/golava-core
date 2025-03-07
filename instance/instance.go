package instance

import (
	"errors"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/auth/generic"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava-core/lang"
	"github.com/wolftotem4/golava-core/routing"
	"github.com/wolftotem4/golava-core/session"
)

type Instance struct {
	App        golava.GolavaApp
	Cookie     cookie.IEncryptableCookieManager
	Session    *session.SessionManager
	Auth       auth.Guard
	Redirector *routing.Redirector
	Locale     string
}

func NewInstance(app golava.GolavaApp) gin.HandlerFunc {
	return func(c *gin.Context) {
		i := &Instance{
			App: app,
			Redirector: &routing.Redirector{
				Router: app.Base().Router,
				GIN:    c,
			},
			Auth: &generic.NullGuard{},
		}

		c.Set("instance", i)

		c.Next()

		i.Dispose()
	}
}

func GetInstance(c *gin.Context) (*Instance, error) {
	obj, ok := c.Get("instance")
	if !ok {
		return nil, errors.New("instance not found in context")
	}
	instance, ok := obj.(*Instance)
	if !ok {
		return nil, errors.New("instance is not of type *Instance")
	}
	return instance, nil
}

func MustGetInstance(c *gin.Context) *Instance {
	instance, err := GetInstance(c)
	if err != nil {
		panic(err)
	}
	return instance
}

func (i *Instance) GetUserPreferredLocale() string {
	if i.Locale == "" {
		return i.App.Base().AppLocale
	}
	return i.Locale
}

func (i *Instance) GetUserPreferredTranslator(options ...lang.TranslatorOption) ut.Translator {
	trans, _ := i.App.Base().Translation.GetTranslator(i.GetUserPreferredLocale())

	args := lang.TranslatorArgs{}
	for _, opt := range options {
		opt(&args)
	}

	return args.Apply(trans)
}

func (i *Instance) Dispose() {
	i.Redirector.GIN = nil
	i.Redirector = nil
	if i.Cookie != nil {
		i.Cookie.SetRequest(nil)
		i.Cookie.SetResponseWriter(nil)
		i.Cookie = nil
	}
	i.App = nil
	i.Session = nil
	i.Auth = nil
}
