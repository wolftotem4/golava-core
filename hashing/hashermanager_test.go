package hashing

import "testing"

func TestIdentifyHasher(t *testing.T) {
	m := NewHasherManager()

	{
		hasher, ok := m.IdentifyHasher("$argon2id$v=19$m=16,t=2,p=1$YTRZaXdqMk11Sms2Q0JQVA$J1Gjx8w3gE4nUxnpneoskA")
		if !ok {
			t.Fatal("hasher not identified")
		}

		if hasher != "argon2id" {
			t.Fatalf("expected argon2id, got %s", hasher)
		}
	}

	{
		hasher, ok := m.IdentifyHasher("$2a$12$Ptw9MMriOubANO6wRQR.quFZs0iD7yBDbONrTMJwB4p3s60oTlqFe")
		if !ok {
			t.Fatal("hasher not identified")
		}

		if hasher != "bcrypt" {
			t.Fatalf("expected bcrypt, got %s", hasher)
		}
	}
}
