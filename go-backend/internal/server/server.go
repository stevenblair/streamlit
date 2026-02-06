package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/stevenblair/streamlit/go-backend/internal/runtime"
)

// Config holds the server configuration
type Config struct {
	Host      string
	Port      int
	BaseURL   string
	ServerURL string
}

// Server represents the Streamlit server
type Server struct {
	config     *Config
	httpServer *http.Server
	runtime    *runtime.Runtime
}

// New creates a new Streamlit server
func New(config *Config) (*Server, error) {
	// Find frontend build directory
	frontendPath, err := findFrontendBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to find frontend build: %w", err)
	}

	// Create runtime
	rt := runtime.New()

	// Create server
	srv := &Server{
		config:  config,
		runtime: rt,
	}

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", srv.handleHealth)

	// WebSocket endpoint for browser connections
	mux.HandleFunc("/_stcore/stream", srv.handleWebSocket)

	// Host config endpoint
	mux.HandleFunc("/_stcore/host-config", srv.handleHostConfig)

	// Media file endpoints
	mux.HandleFunc("/_stcore/media/", srv.handleMedia)

	// Static file server for frontend assets
	fs := http.FileServer(http.Dir(frontendPath))
	mux.Handle("/", fs)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	srv.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv, nil
}

// Start starts the server
func (s *Server) Start() error {
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("\nShutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	// Stop runtime
	s.runtime.Stop()

	log.Println("Server stopped")
	return nil
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// handleHostConfig handles host configuration requests
func (s *Server) handleHostConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	// Return basic host config
	config := `{
		"allowedOrigins": ["*"],
		"useExternalAuthToken": false,
		"enableXsrfProtection": false,
		"enableWebsocketCompression": false,
		"mapboxToken": ""
	}`
	w.Write([]byte(config))
}

// handleMedia handles media file requests
func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement media file handling
	http.NotFound(w, r)
}

// findFrontendBuild locates the frontend build directory
func findFrontendBuild() (string, error) {
	// Get the current executable directory
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exe)

	// Try several possible locations
	candidates := []string{
		// Relative to go-backend during development
		filepath.Join(exeDir, "..", "..", "frontend", "build"),
		filepath.Join(exeDir, "..", "frontend", "build"),
		// Relative to executable
		filepath.Join(exeDir, "frontend", "build"),
		filepath.Join(exeDir, "build"),
		// Hardcoded fallback
		"../frontend/build",
	}

	for _, candidate := range candidates {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		// Check if index.html exists in this directory
		indexPath := filepath.Join(absPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("frontend build directory not found. Please run 'make frontend-fast' from the repo root")
}
