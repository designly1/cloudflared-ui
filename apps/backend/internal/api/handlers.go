package api

import (
	"backend/internal/config"
	"backend/internal/systemd"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from localhost and frontend dev server
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173" || 
		       origin == "http://localhost:3000" || 
		       origin == "http://127.0.0.1:5173" ||
		       origin == "http://127.0.0.1:3000" ||
		       origin == ""
	},
}

type Handlers struct {
	svc *systemd.SystemdService
}

func NewHandlers(svc *systemd.SystemdService) *Handlers {
	return &Handlers{svc: svc}
}

// Response structures
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// StartService handles POST /api/service/start
func (h *Handlers) StartService(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Start(); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service started successfully",
	})
}

// StopService handles POST /api/service/stop
func (h *Handlers) StopService(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Stop(); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service stopped successfully",
	})
}

// RestartService handles POST /api/service/restart
func (h *Handlers) RestartService(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Restart(); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service restarted successfully",
	})
}

// GetStatus handles GET /api/service/status
func (h *Handlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.svc.Status()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    status,
	})
}

// StreamLogs handles WebSocket connection at /api/service/logs
func (h *Handlers) StreamLogs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Create context for streaming
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a custom writer that sends messages over WebSocket
	logWriter := &WebSocketWriter{conn: conn}

	// Handle client disconnection
	conn.SetCloseHandler(func(code int, text string) error {
		cancel()
		return nil
	})

	// Start streaming logs
	if err := h.svc.StreamLogs(ctx, logWriter, true); err != nil {
		log.Printf("Error streaming logs: %v", err)
	}
}

// GetRecentLogs handles GET /api/service/logs/recent
func (h *Handlers) GetRecentLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.svc.GetRecentLogs(100)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    logs,
	})
}

// GetConfig handles GET /api/config
func (h *Handlers) GetConfig(w http.ResponseWriter, r *http.Request) {
	configPath := r.URL.Query().Get("path")
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    cfg,
	})
}

// UpdateConfig handles POST /api/config
func (h *Handlers) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var cfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := config.ValidateConfig(&cfg); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	configPath := r.URL.Query().Get("path")
	if err := config.WriteConfig(configPath, &cfg); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Configuration updated successfully",
	})
}

// HealthCheck handles GET /health
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service is healthy",
	})
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// WebSocketWriter implements io.Writer for WebSocket connections
type WebSocketWriter struct {
	conn *websocket.Conn
}

func (w *WebSocketWriter) Write(p []byte) (n int, err error) {
	if err := w.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return 0, err
	}
	
	if err := w.conn.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	
	return len(p), nil
}

