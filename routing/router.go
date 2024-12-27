package routing

import (
	"net/url"
	"strings"

	"github.com/wolftotem4/golava-core/util"
)

type Router struct {
	BaseURL *url.URL
}

func NewRouter(baseURL string) (*Router, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Router{
		BaseURL: u,
	}, nil
}

func (r *Router) URL(path string) *url.URL {
	u, _ := url.Parse(path)
	return r.BaseURL.ResolveReference(u)
}

// returns the relative path of the given path
func (r *Router) RelativePath(path string) (string, bool) {
	basePath := util.CleanPath(r.BaseURL.Path)
	path = util.CleanPath(path)

	if basePath == path {
		return "", true
	} else if basePath == "" {
		return path, true
	} else if !strings.HasPrefix(path, basePath) {
		return path, false
	}

	value := strings.TrimPrefix(path, basePath)
	if value[0] != '/' {
		return path, false
	}

	return value[1:], true
}
