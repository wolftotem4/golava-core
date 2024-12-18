package hashing

import (
	"strings"

	"github.com/alexedwards/argon2id"
)

var DefaultArgon2idHasher = &Argon2idHasher{Params: argon2id.DefaultParams}

type Argon2idHasher struct {
	Params *argon2id.Params
}

func (h *Argon2idHasher) Make(value string) (string, error) {
	hash, err := argon2id.CreateHash(value, h.Params)
	return hash, err
}

func (h *Argon2idHasher) Check(value string, hashedValue string) (bool, error) {
	return argon2id.ComparePasswordAndHash(value, hashedValue)
}

func (h *Argon2idHasher) NeedsRehash(hashedValue string) bool {
	return !strings.HasPrefix(hashedValue, "$argon2id$")
}
