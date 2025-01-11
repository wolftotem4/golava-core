package auth

import (
	"fmt"
	"strconv"
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
	parts := strings.SplitN(string(r), "|", 3)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func (r Recaller) Hash() string {
	parts := strings.SplitN(string(r), "|", 3)
	if len(parts) > 2 {
		return parts[2]
	}
	return ""
}

func (r Recaller) MatchID(id any) bool {
	switch id.(type) {
	case string:
		return r.ID() == id
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return r.ID() == fmt.Sprintf("%d", id)
	default:
		return false
	}
}

func NewRecaller(id any, token string, hash string) Recaller {
	return Recaller(NewRecallerString(id, token, hash))
}

func NewRecallerString(id any, token string, hash string) string {
	return fmt.Sprintf("%v|%s|%s", id, token, hash)
}

type RecallerIdMorph func(id string) (any, error)

var (
	StringId RecallerIdMorph = func(id string) (any, error) {
		return id, nil
	}

	IntId RecallerIdMorph = func(id string) (any, error) {
		return strconv.Atoi(id)
	}

	Int32Id RecallerIdMorph = func(id string) (any, error) {
		value, err := strconv.ParseInt(id, 10, 32)
		return int32(value), err
	}

	Int64Id RecallerIdMorph = func(id string) (any, error) {
		return strconv.ParseInt(id, 10, 64)
	}
)
