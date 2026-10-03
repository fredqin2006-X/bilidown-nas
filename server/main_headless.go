//go:build headless

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"bilidown/router"
	"bilidown/util"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "check local service health")
	flag.Parse()
	address := os.Getenv("BILIDOWN_LISTEN")
	if address == "" {
		address = ":8098"
	}
	if *healthcheck {
		// Container health probes always connect through loopback.
		port := os.Getenv("BILIDOWN_HEALTH_PORT")
		if port == "" {
			port = "8098"
		}
		client := &http.Client{Timeout: 3 * time.Second}
		response, err := client.Get("http://127.0.0.1:" + port + "/healthz")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}
	os.Setenv("BILIDOWN_SERVER", "1")
	if _, err := util.GetFFmpegPath(); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(util.DatabasePath()), 0755); err != nil {
		log.Fatal(err)
	}
	mustInitTables()
	staticDir := os.Getenv("BILIDOWN_STATIC_DIR")
	if staticDir == "" {
		staticDir = "static"
	}
	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	mux.Handle("/api/", http.StripPrefix("/api", router.API()))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		db := util.MustGetDB()
		defer db.Close()
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "Database unavailable", 503)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
	})
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		log.Print("Stopping HTTP service; unfinished downloads require resubmission after restart")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			server.Close()
		}
	}()
	log.Printf("Bilidown NAS service listening on %s; downloads=%s", address, util.DownloadRoot())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
