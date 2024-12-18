package session

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestInputToMap(t *testing.T) {
	type User struct {
		ID    int
		Name  string
		Email string
	}

	u := User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	m, err := inputToMap(u)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, map[string]interface{}{
		"ID":    1,
		"Name":  "John Doe",
		"Email": "johndoe@example.com",
	}, m)

	m2, err := inputToMap(m)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range m {
		assert.Equal(t, v, m2[k])
	}
}
