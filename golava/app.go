package golava

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/routing"
	"github.com/wolftotem4/golava-core/session"
)

type GolavaApp interface {
	Base() *App
}

type App struct {
	Name           string
	Debug          bool
	AppKey         []byte
	Router         *routing.Router
	Cookie         cookie.IEncryptableCookieManager
	Encryption     encryption.IEncrypter
	Hashing        hashing.Hasher
	SessionFactory *session.SessionFactory
	Translation    *ut.UniversalTranslator
	AppLocale      string
}

func (a *App) Base() *App {
	return a
}
