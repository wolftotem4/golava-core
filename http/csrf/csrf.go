package csrf

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/instance"
)

var ErrTokenMismatch = errors.New("CSRF token mismatch")

func GetCsrfToken(c *gin.Context) string {
	token := c.PostForm("_token")
	if token == "" {
		token = c.GetHeader("X-CSRF-TOKEN")
	}
	if token == "" {
		token = GetCsrfTokenFromXsrf(c)
	}
	return token
}

func GetCsrfTokenFromXsrf(c *gin.Context) string {
	instance := instance.MustGetInstance(c)

	header := c.GetHeader("X-XSRF-TOKEN")
	if header == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return ""
	}

	token, err := instance.App.Base().Encryption.Decrypt(decoded)
	if err != nil {
		return ""
	}

	return string(token)
}

func VerifyCsrfToken(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	if !utils.IsReading(c.Request.Method) {
		if GetCsrfToken(c) != instance.Session.Store.Token() {
			c.Error(ErrTokenMismatch)
			c.Abort()
			return
		}
	}

	addCookieToResponse(instance.App.Base().Cookie, instance.Session.Store.Token(), instance.Session.Lifetime)

	c.Next()
}

func addCookieToResponse(cm cookie.IEncryptableCookieManager, token string, lifetime time.Duration) {
	cm.Encryption().Set(
		"XSRF-TOKEN",
		token,
		cookie.WithMaxAge(int(lifetime.Seconds())),
		cookie.WithHttpOnly(false),
	)
}
