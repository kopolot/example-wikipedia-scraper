package logger

import (
	"log/slog"
	"testing"
)

func TestLogInfo(t *testing.T) {
	Init("test", LevelInfo, true)
	slog.Info("Test log info")
}

func TestLogInScraperFormat(t *testing.T) {
	Init("test", LevelInfo, true)
	Log(LevelInfo, "Test log info with format", "sitename", "example")
}

func TestLogError(t *testing.T) {
	Init("test", LevelError, true)
	Log(LevelDebug, "Test log error", "error", "something went wrong")
	Init("test", LevelDebug, true)
	Log(LevelError, "Test log error with format", "error", "something went wrong")
}
