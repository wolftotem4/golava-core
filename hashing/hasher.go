package hashing

type Hasher interface {
	Make(value string) (string, error)
	Check(value string, hashedValue string) (bool, error)
	NeedsRehash(hashedValue string) bool
}
