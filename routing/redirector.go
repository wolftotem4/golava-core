package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/session"
)

type Redirector struct {
	Router  *Router
	Session *session.SessionManager
	GIN     *gin.Context
}

func (r *Redirector) Redirect(code int, path string) {
	r.GIN.Redirect(code, r.Router.URL(path).String())
}

func (r *Redirector) Intended(code int, path string) {
	r.Session.Store.Flash("url.intended", path)
	r.Redirect(code, path)
}

func (r *Redirector) SetIntendedUrl(url string) {
	r.Session.Store.Flash("url.intended", url)
}

func (r *Redirector) Guest(code int, path string) {
	if r.GIN.Request.Method == "GET" && !utils.ExpectJson(r.GIN.GetHeader("Accept")) {
		r.SetIntendedUrl(r.GIN.Request.URL.String())
	} else {
		r.SetIntendedUrl(r.Previous())
	}

	r.GIN.Redirect(code, r.Router.URL(path).String())
}

func (r *Redirector) Previous(fallback ...string) string {
	referer := r.GIN.Request.Referer()
	if referer == "" {
		if len(fallback) > 0 {
			return r.Router.URL(fallback[0]).String()
		}

		return r.Router.URL("/").String()
	}
	return referer
}

func (r *Redirector) Back(code int, fallback ...string) {
	r.GIN.Redirect(code, r.Previous(fallback...))
}
