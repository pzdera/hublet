package api

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"

	"github.com/pzdera/hublet-v2/internal/config"
)

type Options struct {
	Store *config.Store
	WebFS fs.FS
}

type Server struct {
	store *config.Store
	webFS fs.FS
	mux   *http.ServeMux
}

func New(options Options) (http.Handler, error) {
	if options.Store == nil {
		return nil, errors.New("config store is required")
	}

	if options.WebFS == nil {
		return nil, errors.New("frontend filesystem is required")
	}

	server := &Server{
		store: options.Store,
		webFS: options.WebFS,
		mux:   http.NewServeMux(),
	}

	server.routes()

	return logging(securityHeaders(server.mux)), nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/v2/health", s.health)
	s.mux.HandleFunc("GET /api/v2/config", s.getConfig)
	s.mux.HandleFunc("PUT /api/v2/config", s.putConfig)

	s.mux.Handle("/", http.FileServer(http.FS(s.webFS)))
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *Server) getConfig(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.store.Get())
}

func (s *Server) putConfig(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var cfg config.Config

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.store.Save(cfg); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "same-origin")

		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
