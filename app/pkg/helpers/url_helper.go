package helpers

import "net/url"

func UrlsEqualNormalized(url1 string, url2 string) bool {
	u1, _ := url.Parse(url1)
	u2, _ := url.Parse(url2)

	return u1.Scheme == u2.Scheme &&
		u1.Host == u2.Host &&
		NormalizePath(u1.Path) == NormalizePath(u2.Path) &&
		u1.Query().Encode() == u2.Query().Encode()
}

func NormalizePath(path string) string {
	if len(path) > 0 && path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	return path
}
