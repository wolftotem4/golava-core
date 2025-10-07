package routing

import (
	"testing"
)

func TestURL(t *testing.T) {
	router, err := NewRouter("http://example.com/admin")
	if err != nil {
		t.Error(err)
	}

	url, err := router.URL("home")
	if err != nil {
		t.Error(err)
	}
	if url.String() != "http://example.com/admin/home" {
		t.Errorf("expected http://example.com/admin/home but got %s", url.String())
	}

	url, err = router.URL("/home")
	if err != nil {
		t.Error(err)
	}
	if url.String() != "http://example.com/admin/home" {
		t.Errorf("expected http://example.com/admin/home but got %s", url.String())
	}

	url, err = router.URL("http://example.org/home")
	if err != nil {
		t.Error(err)
	}
	if url.String() != "http://example.org/home" {
		t.Errorf("expected http://example.org/home but got %s", url.String())
	}

	// Test empty path
	url, err = router.URL("")
	if err != nil {
		t.Error(err)
	}
	if url.String() != "http://example.com/admin" {
		t.Errorf("expected http://example.com/admin but got %s", url.String())
	}
}

func TestRelativePath(t *testing.T) {
	router, err := NewRouter("/")
	if err != nil {
		t.Error(err)
	}

	path := "home"
	relativePath, ok := router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != path {
		t.Errorf("expected %s but got %s", path, relativePath)
	}

	path = "/home"
	relativePath, ok = router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != "home" {
		t.Errorf("expected home but got %s", relativePath)
	}

	path = "home"
	router, err = NewRouter("admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != "admin/home" {
		t.Errorf("expected admin/home but got %s", relativePath)
	}

	path = "/home"
	router, err = NewRouter("admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != "admin/home" {
		t.Errorf("expected admin/home but got %s", relativePath)
	}

	path = "http://example.com/admin/home"
	router, err = NewRouter("http://example.com/admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != "home" {
		t.Errorf("expected admin/home but got %s", relativePath)
	}

	path = "http://example.com/home"
	router, err = NewRouter("http://example.com/admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if ok {
		t.Errorf("expected %s to not be relative path", path)
	}

	if relativePath != path {
		t.Errorf("expected %s but got %s", path, relativePath)
	}

	// different host
	path = "http://example.org/admin/home"
	router, err = NewRouter("http://example.com/admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if ok {
		t.Errorf("expected %s to not be relative path", path)
	}

	if relativePath != path {
		t.Errorf("expected %s but got %s", path, relativePath)
	}
}

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		expectErr bool
	}{
		{"valid http URL", "http://example.com/admin", false},
		{"valid https URL", "https://example.com/admin", false},
		{"valid relative path", "/admin", false},
		{"empty string", "", true},
		{"invalid URL", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := NewRouter(tt.baseURL)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if router != nil {
					t.Errorf("expected nil router on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if router == nil {
					t.Errorf("expected valid router")
				}
			}
		})
	}
}

func TestURL_ErrorCases(t *testing.T) {
	router, err := NewRouter("http://example.com/admin")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{"valid relative path", "home", false},
		{"valid absolute path", "/home", false},
		{"valid absolute URL", "http://example.org/home", false},
		{"empty path", "", false},
		{"invalid characters", "home\x00", true}, // null character should cause parsing error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := router.URL(tt.path)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if url == nil {
					t.Errorf("expected valid URL")
				}
			}
		})
	}
}

func TestURLMustPanic(t *testing.T) {
	router, err := NewRouter("http://example.com/admin")
	if err != nil {
		t.Fatal(err)
	}

	// Test normal case
	url := router.URLMustPanic("home")
	if url.String() != "http://example.com/admin/home" {
		t.Errorf("expected http://example.com/admin/home but got %s", url.String())
	}

	// Test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but got none")
		}
	}()
	router.URLMustPanic("home\x00") // Should panic due to invalid character
}

func TestRelativePath_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		baseURL       string
		inputPath     string
		expectedPath  string
		expectedIsRel bool
	}{
		{
			name:          "empty path",
			baseURL:       "http://example.com/admin",
			inputPath:     "",
			expectedPath:  "",
			expectedIsRel: true,
		},
		{
			name:          "exact match",
			baseURL:       "http://example.com/admin",
			inputPath:     "http://example.com/admin",
			expectedPath:  "",
			expectedIsRel: true,
		},
		{
			name:          "different scheme",
			baseURL:       "http://example.com/admin",
			inputPath:     "https://example.com/admin/home",
			expectedPath:  "https://example.com/admin/home",
			expectedIsRel: false,
		},
		{
			name:          "different port",
			baseURL:       "http://example.com:8080/admin",
			inputPath:     "http://example.com:3000/admin/home",
			expectedPath:  "http://example.com:3000/admin/home",
			expectedIsRel: false,
		},
		{
			name:          "root base path",
			baseURL:       "http://example.com/",
			inputPath:     "http://example.com/admin/home",
			expectedPath:  "admin/home",
			expectedIsRel: true,
		},
		{
			name:          "invalid URL in path",
			baseURL:       "/admin",
			inputPath:     "://invalid",
			expectedPath:  "admin/:/invalid",
			expectedIsRel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := NewRouter(tt.baseURL)
			if err != nil {
				t.Fatal(err)
			}

			path, isRel := router.RelativePath(tt.inputPath)
			if path != tt.expectedPath {
				t.Errorf("expected path %q but got %q", tt.expectedPath, path)
			}
			if isRel != tt.expectedIsRel {
				t.Errorf("expected isRelative %v but got %v", tt.expectedIsRel, isRel)
			}
		})
	}
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"/", ""},
		{"admin", "admin"},
		{"/admin", "admin"},
		{"admin/", "admin"},
		{"/admin/", "admin"},
		{"admin/home", "admin/home"},
		{"/admin/home/", "admin/home"},
		{"admin\\home", "admin/home"}, // Test backslash conversion
		{"\\/admin\\home\\/", "admin/home"},
	}

	for _, tt := range tests {
		t.Run("cleanPath_"+tt.input, func(t *testing.T) {
			result := cleanPath(tt.input)
			if result != tt.expected {
				t.Errorf("cleanPath(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
