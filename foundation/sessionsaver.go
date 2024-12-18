package foundation

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
)

func SaveSession(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	c.Next()

	if instance.Session != nil {
		err := instance.Session.Store.Save(c)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Save session error %s", err.Error()))
		}
	}
}
