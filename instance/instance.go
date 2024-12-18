package instance

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava-core/router"
	"github.com/wolftotem4/golava-core/session"
)

type Instance struct {
	App        golava.GolavaApp
	Session    *session.SessionManager
	Auth       auth.Guard
	Redirector *router.Redirector
}

func NewInstance(app golava.GolavaApp) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("instance", &Instance{
			App: app,
			Redirector: &router.Redirector{
				Router: app.Base().Router,
				GIN:    c,
			},
		})

		c.Next()
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
