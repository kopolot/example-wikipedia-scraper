package browser

import "testing"

func TestNormalizeProxyURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "whitespace", in: "  ", want: ""},
		{name: "host port", in: "127.0.0.1:8888", want: "http://127.0.0.1:8888"},
		{name: "http", in: "http://127.0.0.1:8888", want: "http://127.0.0.1:8888"},
		{name: "socks5", in: "socks5://127.0.0.1:1080", want: "socks5://127.0.0.1:1080"},
		{name: "socks4", in: "socks4://proxy.local:1080", want: "socks4://proxy.local:1080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeProxyURL(tt.in); got != tt.want {
				t.Fatalf("NormalizeProxyURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
