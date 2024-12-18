package cookie

import (
	"encoding/base64"
	"net/http"

	"github.com/wolftotem4/golava-core/encryption"
)

type IEncryptableCookieManager interface {
	ICookieManager
	Encryption() ICookieManager
}

type EncryptableCookieManager struct {
	*EncryptCookieManager
}

func NewEncryptableCookieManager(base ICookieManager, encrypter encryption.IEncrypter) *EncryptableCookieManager {
	return &EncryptableCookieManager{
		EncryptCookieManager: &EncryptCookieManager{
			Base:      base,
			Encrypter: encrypter,
		},
	}
}

func (ecm *EncryptableCookieManager) Encryption() ICookieManager {
	return ecm.EncryptCookieManager
}

func (ecm *EncryptableCookieManager) Set(name, value string, options ...WriteOption) {
	ecm.Base.Set(name, value, options...)
}

func (ecm *EncryptableCookieManager) Get(name string) (string, error) {
	return ecm.Base.Get(name)
}

func (ecm *EncryptableCookieManager) Write(cookie *http.Cookie) {
	ecm.Base.Write(cookie)
}

func (ecm *EncryptableCookieManager) Read(name string) (*http.Cookie, error) {
	return ecm.Base.Read(name)
}

func (ecm *EncryptableCookieManager) NewCookie(name, value string) *http.Cookie {
	return ecm.Base.NewCookie(name, value)
}

func (ecm *EncryptableCookieManager) SetRequest(request *http.Request) {
	ecm.Base.SetRequest(request)
}

func (ecm *EncryptableCookieManager) SetResponseWriter(responseWriter http.ResponseWriter) {
	ecm.Base.SetResponseWriter(responseWriter)
}

func (ecm *EncryptableCookieManager) Forget(name string, options ...WriteOption) {
	ecm.Base.Forget(name, options...)
}

type EncryptCookieManager struct {
	Base      ICookieManager
	Encrypter encryption.IEncrypter
}

func (ecm *EncryptCookieManager) Set(name, value string, options ...WriteOption) {
	encrypted, err := ecm.Encrypter.Encrypt([]byte(value))
	if err != nil {
		panic(err)
	}

	ecm.Base.Set(name, base64.StdEncoding.EncodeToString(encrypted), options...)
}

func (ecm *EncryptCookieManager) Get(name string) (string, error) {
	cookie, err := ecm.Base.Get(name)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(cookie)
	if err != nil {
		return "", err
	}

	decrypted, err := ecm.Encrypter.Decrypt(decoded)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (ecm *EncryptCookieManager) Write(cookie *http.Cookie) {
	encrypted, err := ecm.Encrypter.Encrypt([]byte(cookie.Value))
	if err != nil {
		panic(err)
	}

	clonedCookie := *cookie
	clonedCookie.Value = base64.StdEncoding.EncodeToString(encrypted)

	ecm.Base.Write(&clonedCookie)
}

func (ecm *EncryptCookieManager) Read(name string) (*http.Cookie, error) {
	cookie, err := ecm.Base.Read(name)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	decrypted, err := ecm.Encrypter.Decrypt(decoded)
	if err != nil {
		return nil, err
	}

	cookie.Value = string(decrypted)

	return cookie, nil
}

func (ecm *EncryptCookieManager) NewCookie(name, value string) *http.Cookie {
	return ecm.Base.NewCookie(name, value)
}

func (ecm *EncryptCookieManager) SetRequest(request *http.Request) {
	ecm.Base.SetRequest(request)
}

func (ecm *EncryptCookieManager) SetResponseWriter(responseWriter http.ResponseWriter) {
	ecm.Base.SetResponseWriter(responseWriter)
}

func (ecm *EncryptCookieManager) Forget(name string, options ...WriteOption) {
	ecm.Base.Forget(name, options...)
}
