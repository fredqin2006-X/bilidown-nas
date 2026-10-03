package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDownloadFile(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "downloads")
	if err := os.Mkdir(root, 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BILIDOWN_DOWNLOAD_DIR", root)
	inside := filepath.Join(root, "中文视频.mp4")
	outside := filepath.Join(base, "secret.db")
	for _, path := range []string{inside, outside} {
		if err := os.WriteFile(path, []byte("content"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := ResolveDownloadFile(inside); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{outside, root, filepath.Join(root, "..", "secret.db"), filepath.Join(root, "missing.mp4")} {
		if _, err := ResolveDownloadFile(path); err == nil {
			t.Errorf("accepted invalid path: %s", path)
		}
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape.mp4")); err == nil {
		if _, err := ResolveDownloadFile(filepath.Join(root, "escape.mp4")); err == nil {
			t.Error("symlink escaped download root")
		}
	} else {
		t.Logf("symlink test unavailable on this host: %v", err)
	}
}

func TestServerDownloadDirectoryOverridesDatabase(t *testing.T) {
	t.Setenv("BILIDOWN_SERVER", "1")
	root := filepath.Join(t.TempDir(), "downloads")
	t.Setenv("BILIDOWN_DOWNLOAD_DIR", root)
	// Server mode must work even before a database exists.
	got, err := GetCurrentFolder(nil)
	if err != nil || got != root {
		t.Fatalf("folder=%q, error=%v", got, err)
	}
}
