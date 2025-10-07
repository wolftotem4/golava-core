package routing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/session"
)

// mockSessionHandler implements session.SessionHandler interface for testing
type mockSessionHandler struct{}

func (m *mockSessionHandler) Read(ctx context.Context, sessionId string) ([]byte, error) {
	return []byte("{}"), nil
}

func (m *mockSessionHandler) Write(ctx context.Context, sessionId string, data session.SessionData) error {
	return nil
}

func (m *mockSessionHandler) GC(ctx context.Context, lifetime time.Duration) (int64, error) {
	return 0, nil
}

func (m *mockSessionHandler) Destroy(ctx context.Context, sessionId string) error {
	return nil
}

func setupRedirector() (*Redirector, *httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)

	router, err := NewRouter("http://example.com/admin")
	if err != nil {
		panic(err)
	}

	sessionStore := session.NewStore("test-session-id", &mockSessionHandler{})
	sessionStore.Attributes = make(map[string]interface{})

	sessionManager := &session.SessionManager{
		Store: sessionStore,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://example.com/admin/test", nil)

	redirector := &Redirector{
		Router:  router,
		Session: sessionManager,
		GIN:     c,
	}

	return redirector, w, c
}

func TestRedirector_Redirect(t *testing.T) {
	redirector, w, _ := setupRedirector()

	redirector.Redirect(http.StatusFound, "home")

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/home"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}
}

func TestRedirector_Redirect_InvalidPath(t *testing.T) {
	redirector, w, _ := setupRedirector()

	// Use an invalid path that should trigger error handling
	invalidPath := "home\x00"
	redirector.Redirect(http.StatusFound, invalidPath)

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	// Should fallback to the original path
	location := w.Header().Get("Location")
	if location != invalidPath {
		t.Errorf("expected location %s but got %s", invalidPath, location)
	}
}

func TestRedirector_Intended(t *testing.T) {
	redirector, w, _ := setupRedirector()

	// Set intended URL
	redirector.Session.Store.Put("url.intended", "profile")

	redirector.Intended(http.StatusFound, "home")

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/profile"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}

	// Check that intended URL was forgotten
	_, exists := redirector.Session.Store.Get("url.intended")
	if exists {
		t.Error("intended URL should have been forgotten")
	}
}

func TestRedirector_Intended_NoIntendedURL(t *testing.T) {
	redirector, w, _ := setupRedirector()

	redirector.Intended(http.StatusFound, "home")

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/home"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}
}

func TestRedirector_SetIntendedUrl(t *testing.T) {
	redirector, _, _ := setupRedirector()

	redirector.SetIntendedUrl("profile")

	value, exists := redirector.Session.Store.Get("url.intended")
	if !exists {
		t.Error("intended URL should have been set")
	}

	if value != "profile" {
		t.Errorf("expected intended URL to be 'profile' but got %v", value)
	}
}

func TestRedirector_Guest_GET_Request(t *testing.T) {
	redirector, w, c := setupRedirector()

	// Set up GET request without expecting JSON
	c.Request.Method = "GET"
	c.Request.Header.Set("Accept", "text/html")

	redirector.Guest(http.StatusFound, "login")

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/login"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}

	// Check that current URL was set as intended
	intendedURL, exists := redirector.Session.Store.Get("url.intended")
	if !exists {
		t.Error("intended URL should have been set")
	}

	expectedIntended := "http://example.com/admin/test"
	if intendedURL != expectedIntended {
		t.Errorf("expected intended URL %s but got %v", expectedIntended, intendedURL)
	}
}

func TestRedirector_Guest_AJAX_Request(t *testing.T) {
	redirector, _, c := setupRedirector()

	// Set up request expecting JSON (AJAX)
	c.Request.Method = "POST"
	c.Request.Header.Set("Accept", "application/json")
	c.Request.Header.Set("Referer", "http://example.com/admin/previous")

	redirector.Guest(http.StatusFound, "login")

	// Should use previous URL as intended (regardless of redirect status)
	intendedURL, exists := redirector.Session.Store.Get("url.intended")
	if !exists {
		t.Error("intended URL should have been set")
	}

	expectedIntended := "http://example.com/admin/previous"
	if intendedURL != expectedIntended {
		t.Errorf("expected intended URL %s but got %v", expectedIntended, intendedURL)
	}
}

func TestRedirector_Previous(t *testing.T) {
	redirector, _, c := setupRedirector()

	// Test with referer
	c.Request.Header.Set("Referer", "http://example.com/admin/previous")

	result := redirector.Previous()
	expected := "http://example.com/admin/previous"
	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}

	// Test without referer, with fallback
	c.Request.Header.Del("Referer")
	result = redirector.Previous("home")
	expected = "http://example.com/admin/home"
	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}

	// Test without referer, without fallback
	result = redirector.Previous()
	expected = "http://example.com/admin"
	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

func TestRedirector_Previous_URLError(t *testing.T) {
	// Create a router that will cause URL construction to fail
	router, err := NewRouter("http://example.com/admin")
	if err != nil {
		t.Fatal(err)
	}

	sessionStore := session.NewStore("test-session-id", &mockSessionHandler{})
	sessionStore.Attributes = make(map[string]interface{})

	sessionManager := &session.SessionManager{
		Store: sessionStore,
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://example.com/admin/test", nil)

	redirector := &Redirector{
		Router:  router,
		Session: sessionManager,
		GIN:     c,
	}

	// Test fallback when URL construction fails
	result := redirector.Previous("invalid\x00path")
	expected := "invalid\x00path" // Should fallback to raw path
	if result != expected {
		t.Errorf("expected %s but got %s", expected, result)
	}
}

func TestRedirector_Back(t *testing.T) {
	redirector, w, c := setupRedirector()

	c.Request.Header.Set("Referer", "http://example.com/admin/previous")

	redirector.Back(http.StatusFound)

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/previous"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}
}

func TestRedirector_Back_WithFallback(t *testing.T) {
	redirector, w, c := setupRedirector()

	// No referer header
	c.Request.Header.Del("Referer")

	redirector.Back(http.StatusFound, "home")

	if w.Code != http.StatusFound {
		t.Errorf("expected status %d but got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	expected := "http://example.com/admin/home"
	if location != expected {
		t.Errorf("expected location %s but got %s", expected, location)
	}
}
