package auth

import (
	"fmt"
	"strings"
)

type Recaller string

func (r Recaller) Valid() bool {
	parts := strings.SplitN(string(r), "|", 3)
	return len(parts) == 3 || strings.TrimSpace(parts[0]) != "" || strings.TrimSpace(parts[1]) != ""
}

func (r Recaller) ID() string {
	return strings.SplitN(string(r), "|", 3)[0]
}

func (r Recaller) Token() string {
	return strings.SplitN(string(r), "|", 3)[1]
}

func (r Recaller) Hash() string {
	return strings.SplitN(string(r), "|", 3)[2]
}

func NewRecaller(id any, token string, hash string) Recaller {
	return Recaller(NewRecallerString(id, token, hash))
}

func NewRecallerString(id any, token string, hash string) string {
	return fmt.Sprintf("%v|%s|%s", id, token, hash)
}
