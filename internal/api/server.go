package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the HTTP API server
type Server struct {
	router *mux.Router
	store  Store
	server *http.Server
}

// Store defines the interface for the key-value store
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// NewServer creates a new HTTP API server
func NewServer(store Store, addr string) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.setupRoutes()

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/v1/kv/{key}", s.handleGet).Methods("GET")
	s.router.HandleFunc("/v1/kv/{key}", s.handleSet).Methods("PUT")
	s.router.HandleFunc("/v1/kv/{key}", s.handleDelete).Methods("DELETE")
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := s.store.Get(r.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key":   key,
		"value": value,
	})
}

type setRequest struct {
	Value string `json:"value"`
	TTL   string `json:"ttl,omitempty"`
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var req setRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ttl time.Duration
	if req.TTL != "" {
		var err error
		ttl, err = time.ParseDuration(req.TTL)
		if err != nil {
			http.Error(w, "invalid TTL format", http.StatusBadRequest)
			return
		}
	}

	if err := s.store.Set(r.Context(), key, req.Value, ttl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := s.store.Delete(r.Context(), key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
