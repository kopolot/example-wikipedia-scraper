package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

var rootRepoPath string

func FindRepoRoot(startDir string) string {
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

func GetAbsoluteFilePath(relativeToCurrentFilePath string) string {
	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	_, currentFilePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	thisFileDir := filepath.Dir(thisFilePath)
	projetcRootDir := filepath.Dir(filepath.Dir(thisFileDir))
	currentFileDir := filepath.Dir(currentFilePath)
	relDir, err := filepath.Rel(projetcRootDir, currentFileDir)
	if err != nil {
		return ""
	}
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	execPathDir := filepath.Dir(execPath)
	absPath := filepath.Join(execPathDir, relDir, relativeToCurrentFilePath)
	if _, fileExists := os.Stat(absPath); fileExists == nil {
		return absPath
	}
	rootRepoPath = FindRepoRoot(execPathDir)
	if rootRepoPath == "" {
		rootRepoPath = projetcRootDir
	}
	absPath = filepath.Join(rootRepoPath, relDir, relativeToCurrentFilePath)
	return absPath
}

func GetCurrentFilePath() string {
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
	rootRepoPath = FindRepoRoot(execPathDir)
	absPath := filepath.Join(rootRepoPath, relPos)
	return absPath
}
