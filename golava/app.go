package golava

import (
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/router"
	"github.com/wolftotem4/golava-core/session"
)

type GolavaApp interface {
	Base() *App
}

type App struct {
	Name           string
	Debug          bool
	AppKey         []byte
	Router         *router.Router
	Cookie         cookie.IEncryptableCookieManager
	Encryption     encryption.IEncrypter
	Hashing        hashing.Hasher
	SessionFactory *session.SessionFactory
}

func (a *App) Base() *App {
	return a
}
