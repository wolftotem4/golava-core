package routing

import (
	"net/url"
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
