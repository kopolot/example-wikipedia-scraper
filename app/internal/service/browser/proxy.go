package browser

import "strings"

func NormalizeProxyURL(raw string) string {
	proxy := strings.TrimSpace(raw)
	if proxy == "" {
		return ""
	}
	lower := strings.ToLower(proxy)
	if strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "socks4://") ||
		strings.HasPrefix(lower, "socks5://") {
		return proxy
	}
	return "http://" + proxy
}

func cloneEngineSettings(src map[string]any) map[string]any {
	if src == nil {
		return make(map[string]any)
	}
	dst := make(map[string]any, len(src)+1)
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
