package hashing

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

const Md5Prefix = "$md5$"

type Md5Hasher struct {
	Callback func(value string) (string, error)
}

func (h *Md5Hasher) make(value string) (hashedValue [16]byte, err error) {
	if h.Callback != nil {
		value, err = h.Callback(value)
		if err != nil {
			return [16]byte{}, err
		}
	}

	hash := md5.Sum([]byte(value))
	return hash, nil
}

func (h *Md5Hasher) Make(value string) (string, error) {
	hash, err := h.make(value)
	return fmt.Sprintf("%s%x", Md5Prefix, hash), err
}

func (h *Md5Hasher) Check(value string, hashedValue string) (bool, error) {
	otherHash, err := h.make(value)
	if err != nil {
		return false, err
	}

	hashedBytes, err := hex.DecodeString(h.stripPrefix(hashedValue))
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(hashedBytes, otherHash[:]) == 1, nil
}

func (h *Md5Hasher) NeedsRehash(hashedValue string) bool {
	return false
}

func (h *Md5Hasher) stripPrefix(hashedValue string) string {
	return hashedValue[len(Md5Prefix):]
}
