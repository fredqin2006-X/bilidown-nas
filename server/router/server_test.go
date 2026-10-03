package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"bilidown/util"
	_ "modernc.org/sqlite"
)

func TestServerRestrictions(t *testing.T) {
	t.Setenv("BILIDOWN_SERVER", "1")
	t.Setenv("BILIDOWN_DB_PATH", filepath.Join(t.TempDir(), "data.db"))
	handler := API()
	for _, path := range []string{"/quit", "/showFile?filePath=/etc/passwd", "/downloadVideo?path=/etc/passwd"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code < 400 {
			t.Errorf("%s returned %d", path, response.Code)
		}
	}
}

func TestTaskDownloadSupportsRangeAndRejectsOutsidePath(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BILIDOWN_SERVER", "1")
	t.Setenv("BILIDOWN_DOWNLOAD_DIR", root)
	t.Setenv("BILIDOWN_DB_PATH", filepath.Join(t.TempDir(), "data.db"))
	db := util.MustGetDB()
	defer db.Close()
	_, err := db.Exec(`CREATE TABLE task (id INTEGER PRIMARY KEY, bvid TEXT, cid INTEGER, format INTEGER, title TEXT, owner TEXT, cover TEXT, status TEXT, folder TEXT, duration INTEGER, download_type TEXT, create_at TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO task VALUES (1,'BV1',1,80,'测试','owner','','done',?,1,'merge','2026-10-03 00:00:00')`, root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "测试 MQ.mp4"), []byte("0123456789"), 0600); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("GET", "/downloadVideo?id=1", nil)
	request.Header.Set("Range", "bytes=2-5")
	response := httptest.NewRecorder()
	downloadVideo.ServeHTTP(response, request)
	if response.Code != http.StatusPartialContent || response.Body.String() != "2345" {
		t.Fatalf("range response=%d %q", response.Code, response.Body.String())
	}
	if _, err := db.Exec(`UPDATE task SET folder=? WHERE id=1`, filepath.Dir(root)); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	downloadVideo.ServeHTTP(response, httptest.NewRequest("GET", "/downloadVideo?id=1", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("outside path status=%d", response.Code)
	}
}
