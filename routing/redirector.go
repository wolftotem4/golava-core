package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/session"
)

// Redirector provides convenient methods for HTTP redirections in Gin applications
// It integrates with Router for URL generation and SessionManager for session handling
type Redirector struct {
	Router  *Router                 // Router instance for URL generation
	Session *session.SessionManager // Session manager for storing intended URLs
	GIN     *gin.Context            // Gin context for HTTP operations
}

// Redirect performs an HTTP redirect to the specified path with the given status code
// The path is resolved through the Router to generate the full URL
func (r *Redirector) Redirect(code int, path string) {
	url, err := r.Router.URL(path)
	if err != nil {
		// Fallback to the path itself if URL construction fails
		r.GIN.Redirect(code, path)
		return
	}
	r.GIN.Redirect(code, url.String())
}

// Intended redirects to the previously intended URL stored in session, or defaults if none exists
// After redirecting, the intended URL is removed from the session
func (r *Redirector) Intended(code int, defaults string) {
	path, ok := r.Session.Store.Attributes["url.intended"].(string)
	if !ok {
		path = defaults
	}
	r.Session.Store.Forget("url.intended")

	r.Redirect(code, path)
}

// SetIntendedUrl stores the given URL as the intended destination in the session
// This is typically used to remember where a user wanted to go before authentication
func (r *Redirector) SetIntendedUrl(url string) {
	r.Session.Store.Put("url.intended", url)
}

// Guest redirects unauthenticated users to the specified path
// For GET requests (non-AJAX), it stores the current URL as intended destination
// For other requests or AJAX requests, it uses the previous URL as intended destination
func (r *Redirector) Guest(code int, path string) {
	if r.GIN.Request.Method == "GET" && !utils.ExpectJson(r.GIN.GetHeader("Accept")) {
		r.SetIntendedUrl(r.GIN.Request.URL.String())
	} else {
		r.SetIntendedUrl(r.Previous())
	}

	url, err := r.Router.URL(path)
	if err != nil {
		// Fallback to the path itself if URL construction fails
		r.GIN.Redirect(code, path)
		return
	}
	r.GIN.Redirect(code, url.String())
}

// Previous returns the previous URL from the HTTP referer header
// If no referer is present, it returns the fallback URL or root URL
// The fallback path is resolved through the Router if possible
func (r *Redirector) Previous(fallback ...string) string {
	referer := r.GIN.Request.Referer()
	if referer == "" {
		if len(fallback) > 0 {
			if url, err := r.Router.URL(fallback[0]); err == nil {
				return url.String()
			}
			return fallback[0] // Fallback to raw path if URL construction fails
		}

		if url, err := r.Router.URL("/"); err == nil {
			return url.String()
		}
		return "/" // Fallback to root path
	}
	return referer
}

// Back redirects to the previous URL (from referer header) with the given status code
// If no referer is available, it uses the fallback URL or root URL
func (r *Redirector) Back(code int, fallback ...string) {
	r.GIN.Redirect(code, r.Previous(fallback...))
}
