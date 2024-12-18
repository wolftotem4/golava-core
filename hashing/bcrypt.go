package hashing

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var DefaultBcryptHasher = &BcryptHasher{Cost: bcrypt.DefaultCost}

type BcryptHasher struct {
	Cost int
}

func (h *BcryptHasher) Make(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), h.Cost)
	return string(bytes), err
}

func (h *BcryptHasher) Check(value string, hashedValue string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	if err == nil {
		return true, nil
	}

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}

func (h *BcryptHasher) NeedsRehash(hashedValue string) bool {
	if !strings.HasPrefix(hashedValue, "$2a$") && !strings.HasPrefix(hashedValue, "$2b$") {
		return true
	}

	cost, err := bcrypt.Cost([]byte(hashedValue))
	if err != nil {
		return false
	}

	return cost != h.Cost
}
