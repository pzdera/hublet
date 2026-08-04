package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pzdera/hublet/internal/api"
	"github.com/pzdera/hublet/internal/config"
)

//go:embed all:web-dist
var frontend embed.FS

func main() {
	addr := envOr(
		"HUBLET_ADDR",
		":3000",
	)

	dataDir := envOr(
		"HUBLET_DATA_DIR",
		"/app/data",
	)

	wallpaperDir := envOr(
		"HUBLET_WALLPAPER_DIR",
		"/app/wallpapers",
	)

	iconDir := envOr(
		"HUBLET_ICON_DIR",
		"/app/icons",
	)

	absoluteWallpaperDir, err :=
		prepareDirectory(wallpaperDir)

	if err != nil {
		log.Fatalf(
			"initialize wallpaper directory: %v",
			err,
		)
	}

	absoluteIconDir, err :=
		prepareDirectory(iconDir)

	if err != nil {
		log.Fatalf(
			"initialize icon directory: %v",
			err,
		)
	}

	store, err := config.NewStore(dataDir)

	if err != nil {
		log.Fatalf(
			"initialize config store: %v",
			err,
		)
	}

	webFS, err := fs.Sub(
		frontend,
		"web-dist",
	)

	if err != nil {
		log.Fatalf(
			"load frontend: %v",
			err,
		)
	}

	handler, err := api.New(api.Options{
		Store:        store,
		WebFS:        webFS,
		WallpaperDir: absoluteWallpaperDir,
		IconDir:      absoluteIconDir,
	})

	if err != nil {
		log.Fatalf(
			"initialize server: %v",
			err,
		)
	}

	log.Printf(
		"Hublet listening on %s",
		addr,
	)

	log.Printf(
		"Wallpaper directory: %s",
		absoluteWallpaperDir,
	)

	log.Printf(
		"Icon directory: %s",
		absoluteIconDir,
	)

	if err := http.ListenAndServe(
		addr,
		handler,
	); err != nil {
		log.Fatal(err)
	}
}

func prepareDirectory(
	path string,
) (string, error) {
	if err := os.MkdirAll(
		path,
		0o755,
	); err != nil {
		return "", err
	}

	return filepath.Abs(path)
}

func envOr(
	name string,
	fallback string,
) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
