package utils

import (
	"net/http"
	"strings"
)

func ExpectJson(accept string) bool {
	for _, accept := range parseAcceptHeader(accept) {
		if strings.HasPrefix(accept, "application/json") {
			return true
		}
	}
	return false
}

func parseAcceptHeader(accept string) []string {
	accepts := strings.Split(accept, ",")
	for i, accept := range accepts {
		accepts[i] = strings.TrimSpace(accept)
	}
	return accepts
}

func IsReading(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
