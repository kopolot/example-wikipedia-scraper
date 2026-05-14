package helpers

import (
	"path/filepath"
	"runtime"
)

func GetCurrentFilePath() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Dir(file)
}
