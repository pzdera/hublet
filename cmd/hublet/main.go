package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/pzdera/hublet-v2/internal/api"
	"github.com/pzdera/hublet-v2/internal/config"
)

//go:embed all:web-dist
var frontend embed.FS

func main() {
	addr := envOr("HUBLET_ADDR", ":3000")
	dataDir := envOr("HUBLET_DATA_DIR", "/app/data")

	store, err := config.NewStore(dataDir)
	if err != nil {
		log.Fatalf("initialize config store: %v", err)
	}

	webFS, err := fs.Sub(frontend, "web-dist")
	if err != nil {
		log.Fatalf("load frontend: %v", err)
	}

	handler, err := api.New(api.Options{
		Store: store,
		WebFS: webFS,
	})
	if err != nil {
		log.Fatalf("initialize server: %v", err)
	}

	log.Printf("Hublet v2 listening on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOr(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
