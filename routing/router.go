package routing

import (
	"errors"
	"net/url"
	pathLib "path"
	"strings"
)

// ErrInvalidURL is returned when an invalid URL is provided
var ErrInvalidURL = errors.New("invalid URL provided")

// Router provides URL routing functionality with a base URL
type Router struct {
	BaseURL *url.URL
}

// NewRouter creates a new Router instance with the given base URL
// baseURL should be a valid URL string (e.g., "http://example.com/admin" or "/admin")
func NewRouter(baseURL string) (*Router, error) {
	if baseURL == "" {
		return nil, ErrInvalidURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Router{
		BaseURL: u,
	}, nil
}

// URL constructs a full URL by joining the given path with the base URL
// If the path is already an absolute URL, it returns that URL unchanged
// If the path is relative, it joins it with the router's base URL
func (r *Router) URL(path string) (*url.URL, error) {
	if path == "" {
		return r.BaseURL, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if u.IsAbs() {
		return u, nil
	}

	// Clean the path by removing leading slash to ensure proper joining
	cleanedPath := strings.TrimPrefix(path, "/")
	joinedPath := pathLib.Join(r.BaseURL.Path, cleanedPath)

	u, err = url.Parse(joinedPath)
	if err != nil {
		return nil, err
	}

	return r.BaseURL.ResolveReference(u), nil
}

// URLMustPanic is a convenience method that returns the URL or panics on error
// Use this only when you are certain the path is valid
func (r *Router) URLMustPanic(path string) *url.URL {
	u, err := r.URL(path)
	if err != nil {
		panic(err)
	}
	return u
}

// RelativePath returns the relative path of the given path against the router's base URL
// The second return value indicates whether the path is relative to the base URL
// If the path is not relative to the base URL, the original path is returned
func (r *Router) RelativePath(path string) (string, bool) {
	if path == "" {
		return "", true
	}

	u, err := url.Parse(path)
	if err != nil {
		// If parsing fails, treat as relative path
		return cleanPath(pathLib.Join(r.BaseURL.Path, path)), true
	}

	if u.IsAbs() {
		// Different scheme, host, or port means not relative
		if u.Scheme != r.BaseURL.Scheme || u.Host != r.BaseURL.Host {
			return path, false
		}

		basePath := cleanPath(r.BaseURL.Path)
		targetPath := cleanPath(u.Path)

		// Exact match
		if basePath == targetPath {
			return "", true
		}

		// If base path is empty, return the target path
		if basePath == "" {
			return targetPath, true
		}

		// Check if target path starts with base path
		if !strings.HasPrefix(targetPath, basePath) {
			return path, false
		}

		// Remove base path from target path
		relativePath := strings.TrimPrefix(targetPath, basePath)

		// Ensure the remaining path starts with '/' (proper prefix)
		if len(relativePath) == 0 {
			return "", true
		}

		if relativePath[0] != '/' {
			return path, false
		}

		return relativePath[1:], true
	}

	// For relative paths, join with base path
	return cleanPath(pathLib.Join(r.BaseURL.Path, path)), true
}

// cleanPath normalizes a path by converting backslashes to forward slashes
// and removing leading/trailing slashes
func cleanPath(path string) string {
	if path == "" {
		return ""
	}
	// Convert backslashes to forward slashes for cross-platform compatibility
	cleaned := strings.ReplaceAll(path, "\\", "/")
	// Remove leading and trailing slashes
	return strings.Trim(cleaned, "/")
}
