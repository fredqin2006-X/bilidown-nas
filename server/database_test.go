//go:build headless

package main

import (
	"bilidown/util"
	"path/filepath"
	"testing"
)

func TestHeadlessDatabaseStartup(t *testing.T) {
	t.Setenv("BILIDOWN_SERVER", "1")
	t.Setenv("BILIDOWN_DB_PATH", filepath.Join(t.TempDir(), "data.db"))
	t.Setenv("BILIDOWN_DOWNLOAD_DIR", filepath.Join(t.TempDir(), "downloads"))
	mustInitTables()
	db := util.MustGetDB()
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO task (bvid,cid,format,title,owner,cover,status,folder,duration) VALUES ('BV1TEST',1,80,'test','test','','waiting',?,1)`, util.DownloadRoot()); err != nil {
		t.Fatal(err)
	}
	mustInitTables()
	var status string
	if err := db.QueryRow(`SELECT status FROM task`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "error" {
		t.Fatalf("interrupted task status=%s", status)
	}
}
