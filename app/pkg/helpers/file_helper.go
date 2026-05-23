package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

var rootRepoPath string

func GetCurrentFilePath() string {
	// implemented like this because runtime returns path to .go file , not based on currenrt working directory, which can be different when running tests from different locations
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	_, targetFile, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	currentDir := filepath.Dir(currentFile)
	rootDir := filepath.Dir(filepath.Dir(currentDir))
	targetDir := filepath.Dir(targetFile)
	relPos, err := filepath.Rel(rootDir, targetDir)
	if err != nil {
		return ""
	}
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	execPathDir := filepath.Dir(execPath)
	rootRepoPath = findRepoRoot(execPathDir)
	absPath := filepath.Join(rootRepoPath, relPos)
	return absPath
}

func findRepoRoot(startDir string) string {
	if rootRepoPath == "" {
		dir := startDir
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				rootRepoPath = dir
				break
			}
			if parent := filepath.Dir(dir); parent == dir {
				dir = ""
				break
			} else {
				dir = parent
			}
		}
	}
	return rootRepoPath
}
