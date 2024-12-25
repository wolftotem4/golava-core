package utils

import (
	"net/http"
	"strings"
)

func ExpectJson(accepts []string) bool {
	for _, accept := range accepts {
		if strings.HasPrefix(accept, "application/json") {
			return true
		}
	}
	return false
}

func IsReading(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
