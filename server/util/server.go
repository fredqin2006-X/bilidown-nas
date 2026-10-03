package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ServerMode() bool { return os.Getenv("BILIDOWN_SERVER") == "1" }
func DatabasePath() string {
	if value := os.Getenv("BILIDOWN_DB_PATH"); value != "" {
		return value
	}
	return "./data.db"
}
func DownloadRoot() string {
	if value := os.Getenv("BILIDOWN_DOWNLOAD_DIR"); value != "" {
		return value
	}
	return "./download"
}
func within(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}
func WithinDownloadRoot(path string) bool {
	root, err := filepath.Abs(DownloadRoot())
	if err != nil {
		return false
	}
	absolute, err := filepath.Abs(path)
	return err == nil && within(root, absolute)
}

// ResolveDownloadFile checks both lexical paths and symlinks before serving files.
func ResolveDownloadFile(path string) (string, error) {
	if !WithinDownloadRoot(path) {
		return "", errors.New("outside download directory")
	}
	root, err := filepath.EvalSymlinks(DownloadRoot())
	if err != nil {
		return "", err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	resolved, err = filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	if !within(root, resolved) {
		return "", errors.New("outside download directory")
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("not a regular file")
	}
	return resolved, nil
}
