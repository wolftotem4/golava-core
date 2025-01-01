package hashing

import (
	"errors"
	"strings"
)

var ErrUnknownHasher = errors.New("unknown hasher")

var deprecatedHashers = map[string]struct{}{
	"bcrypt": {},
}

type HasherManager struct {
	DefaultHasher string
	Hashers       map[string]Hasher
	MapHashPrefix map[string]string
}

func NewHasherManager() *HasherManager {
	return &HasherManager{
		DefaultHasher: "argon2id",
		Hashers: map[string]Hasher{
			"bcrypt":   DefaultBcryptHasher,
			"argon2id": DefaultArgon2idHasher,
		},
		MapHashPrefix: map[string]string{
			"$2a$":       "bcrypt",
			"$2b$":       "bcrypt",
			"$argon2id$": "argon2id",
		},
	}
}

func (m *HasherManager) Make(value string) (string, error) {
	return m.Hashers[m.DefaultHasher].Make(value)
}

func (m *HasherManager) Check(value string, hashedValue string) (bool, error) {
	hasher, ok := m.IdentifyHasher(hashedValue)
	if !ok {
		return false, ErrUnknownHasher
	}

	return m.Hashers[hasher].Check(value, hashedValue)
}

func (m *HasherManager) NeedsRehash(hashedValue string) bool {
	hasher, _ := m.IdentifyHasher(hashedValue)
	if IsDeprecated(hasher) && hasher != m.DefaultHasher {
		return true
	}

	return m.Hashers[hasher].NeedsRehash(hashedValue)
}

func (m *HasherManager) IdentifyHasher(hashedValue string) (string, bool) {
	segments := strings.SplitN(hashedValue, "$", 3)

	var sb strings.Builder
	for i := 0; i < 2; i++ {
		sb.WriteString(segments[i])
		sb.WriteString("$")
	}
	prefix := sb.String()

	hasher, ok := m.MapHashPrefix[prefix]
	return hasher, ok
}

func IsDeprecated(hasher string) bool {
	_, ok := deprecatedHashers[hasher]
	return ok
}

func MarkDeprecated(hasher string, hashers ...string) {
	deprecatedHashers[hasher] = struct{}{}
	for _, h := range hashers {
		deprecatedHashers[h] = struct{}{}
	}
}

func UnmarkDeprecated(hasher string, hashers ...string) {
	delete(deprecatedHashers, hasher)
	for _, h := range hashers {
		delete(deprecatedHashers, h)
	}
}
