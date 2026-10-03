package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bilidown/task"
	"bilidown/util"
)

func getServerInfo(w http.ResponseWriter, r *http.Request) {
	hostPath := os.Getenv("BILIDOWN_DOWNLOAD_HOST_PATH")
	if hostPath == "" {
		hostPath = util.DownloadRoot()
	}
	util.Res{Success: true, Data: struct {
		ServerMode       bool   `json:"serverMode"`
		DownloadFolder   string `json:"downloadFolder"`
		DownloadHostPath string `json:"downloadHostPath"`
	}{util.ServerMode(), util.DownloadRoot(), hostPath}}.Write(w)
}

var downloadVideo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Query().Get("path")
	if idValue := r.URL.Query().Get("id"); idValue != "" {
		id, err := strconv.Atoi(idValue)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid task ID", 400)
			return
		}
		db := util.MustGetDB()
		defer db.Close()
		item, err := task.GetTask(db, id)
		if err != nil || item.Status != "done" {
			http.NotFound(w, r)
			return
		}
		path = item.FilePath()
	} else if util.ServerMode() {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}
	if util.ServerMode() {
		resolved, err := util.ResolveDownloadFile(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		path = resolved
	} else {
		path = strings.ReplaceAll(filepath.Clean(path), "\\", "/")
	}
	http.ServeFile(w, r, path)
})
