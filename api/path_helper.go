package api

import "strings"

func NormalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func ToWindowsPath(path string) string {
	return strings.ReplaceAll(path, "/", "\\")
}
