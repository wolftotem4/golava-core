package routing

import "testing"

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

	path = "home"
	router, err = NewRouter("admin")
	if err != nil {
		t.Error(err)
	}

	relativePath, ok = router.RelativePath(path)
	if ok {
		t.Errorf("expected %s to be not relative path", path)
	}

	path = "admin/home"
	relativePath, ok = router.RelativePath(path)
	if !ok {
		t.Errorf("expected %s to be relative path", path)
	}

	if relativePath != "home" {
		t.Errorf("expected home but got %s", relativePath)
	}
}
