package util

import "strings"

func CleanPath(path string) string {
	return strings.TrimLeft(strings.ReplaceAll(path, "\\", "/"), "/")
}
