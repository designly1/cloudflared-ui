package api

import (
	"backend/internal/systemd"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates and configures the API router
func NewRouter(svc *systemd.SystemdService) *mux.Router {
	h := NewHandlers(svc)
	r := mux.NewRouter()

	// Enable CORS for development
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Service control endpoints
	api.HandleFunc("/service/start", h.StartService).Methods("POST")
	api.HandleFunc("/service/stop", h.StopService).Methods("POST")
	api.HandleFunc("/service/restart", h.RestartService).Methods("POST")
	api.HandleFunc("/service/status", h.GetStatus).Methods("GET")
	
	// Logs endpoints
	api.HandleFunc("/service/logs", h.StreamLogs)
	api.HandleFunc("/service/logs/recent", h.GetRecentLogs).Methods("GET")
	
	// Config endpoints
	api.HandleFunc("/config", h.GetConfig).Methods("GET")
	api.HandleFunc("/config", h.UpdateConfig).Methods("POST")

	// Health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")

	return r
}

// CORS middleware for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Allow requests from localhost during development
		if origin == "http://localhost:5173" || 
		   origin == "http://localhost:3000" || 
		   origin == "http://127.0.0.1:5173" ||
		   origin == "http://127.0.0.1:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

