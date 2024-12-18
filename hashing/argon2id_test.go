package hashing

import (
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestArgon2idHasher_Make(t *testing.T) {
	h := Argon2idHasher{Params: argon2id.DefaultParams}
	hash, err := h.Make("password")
	if err != nil {
		t.Fatal(err)
	}

	if hash == "" {
		t.Fatal("hash is empty")
	}
}

func TestArgon2idHasher_Check(t *testing.T) {
	h := Argon2idHasher{Params: argon2id.DefaultParams}

	ok, err := h.Check("password", "$argon2id$v=19$m=16,t=2,p=1$YTRZaXdqMk11Sms2Q0JQVA$J1Gjx8w3gE4nUxnpneoskA")
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("passwords do not match")
	}
}
