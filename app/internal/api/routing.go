package api

import "strings"

func joinRoutePath(prefix, path string) string {
	prefix = strings.Trim(prefix, "/")
	path = strings.TrimPrefix(path, "/")

	if prefix == "" {
		if path == "" {
			return "/"
		}
		return "/" + path
	}
	if path == "" {
		return prefix + "/"
	}
	return prefix + "/" + path
}
