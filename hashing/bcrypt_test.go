package hashing

import "testing"

func TestBcryptHasher_Make(t *testing.T) {
	h := BcryptHasher{Cost: 10}
	hash, err := h.Make("password")
	if err != nil {
		t.Fatal(err)
	}

	if hash == "" {
		t.Fatal("hash is empty")
	}
}

func TestBcryptHasher_Check(t *testing.T) {
	h := BcryptHasher{Cost: 12}

	ok, err := h.Check("password", "$2a$12$Ptw9MMriOubANO6wRQR.quFZs0iD7yBDbONrTMJwB4p3s60oTlqFe")
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("passwords do not match")
	}
}
