package cookie

import (
	"net/http"
	"time"
)

type ICookieManager interface {
	Set(name, value string, options ...WriteOption)
	Get(name string) (string, error)
	Write(cookie *http.Cookie)
	Read(name string) (*http.Cookie, error)
	NewCookie(name, value string) *http.Cookie
	SetRequest(request *http.Request)
	SetResponseWriter(responseWriter http.ResponseWriter)
	Forget(name string, options ...WriteOption)
}

type CookieManager struct {
	Path           string
	Domain         string
	Secure         bool
	SameSite       http.SameSite
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func (cm *CookieManager) Set(name, value string, options ...WriteOption) {
	cookie := cm.NewCookie(name, value)

	for _, option := range options {
		option(cookie)
	}

	cm.Write(cookie)
}

func (cm *CookieManager) Get(name string) (string, error) {
	cookie, err := cm.Read(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (cm *CookieManager) Write(cookie *http.Cookie) {
	http.SetCookie(cm.ResponseWriter, cookie)
}

func (cm *CookieManager) Read(name string) (*http.Cookie, error) {
	return cm.Request.Cookie(name)
}

func (cm *CookieManager) NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cm.Path,
		Domain:   cm.Domain,
		Secure:   cm.Secure,
		SameSite: cm.SameSite,
		HttpOnly: true,
	}
}

func (cm *CookieManager) SetRequest(request *http.Request) {
	cm.Request = request
}

func (cm *CookieManager) SetResponseWriter(responseWriter http.ResponseWriter) {
	cm.ResponseWriter = responseWriter
}

func (cm *CookieManager) Forget(name string, options ...WriteOption) {
	cookie := cm.NewCookie(name, "")

	for _, option := range options {
		option(cookie)
	}

	cookie.MaxAge = 0
	cookie.Expires = time.Unix(0, 0)

	cm.Write(cookie)
}

type WriteOption func(*http.Cookie)

func WithMaxAge(maxAge int) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.MaxAge = maxAge
	}
}

func WithExpires(expires time.Time) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.Expires = expires
	}
}

func WithPath(path string) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.Path = path
	}
}

func WithDomain(domain string) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.Domain = domain
	}
}

func WithSecure(secure bool) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.Secure = secure
	}
}

func WithHttpOnly(httpOnly bool) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.HttpOnly = httpOnly
	}
}

func WithSameSite(sameSite http.SameSite) WriteOption {
	return func(cookie *http.Cookie) {
		cookie.SameSite = sameSite
	}
}
