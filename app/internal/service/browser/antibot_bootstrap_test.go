package browser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAbckCookieValueValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		value := stringsRepeat("A", 120) + "~0~0~1234567890"
		assert.True(t, isAbckCookieValueValid(value))
	})

	t.Run("pending challenge", func(t *testing.T) {
		value := stringsRepeat("A", 120) + "~-1~-1~1783345412"
		assert.False(t, isAbckCookieValueValid(value))
	})

	t.Run("too short", func(t *testing.T) {
		assert.False(t, isAbckCookieValueValid("short"))
	})
}

func stringsRepeat(s string, count int) string {
	out := make([]byte, 0, len(s)*count)
	for range count {
		out = append(out, s...)
	}
	return string(out)
}
