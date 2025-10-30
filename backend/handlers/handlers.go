package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/phughe11/catiSip/backend/config"
	"github.com/phughe11/catiSip/backend/sip"
)

// Handler contains dependencies for HTTP handlers
type Handler struct {
	sipClient *sip.Client
	config    *config.Config
}

// New creates a new handler instance
func New(sipClient *sip.Client, cfg *config.Config) *Handler {
	return &Handler{
		sipClient: sipClient,
		config:    cfg,
	}
}

// HealthCheck returns the health status of the service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"service": "catiSip",
	}
	writeJSON(w, http.StatusOK, response)
}

// MakeCall handles call initiation requests
func (h *Handler) MakeCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.From == "" || req.To == "" {
		http.Error(w, "Both 'from' and 'to' fields are required", http.StatusBadRequest)
		return
	}

	call, err := h.sipClient.MakeCall(req.From, req.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, call)
}

// HangupCall handles call termination requests
func (h *Handler) HangupCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID string `json:"call_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" {
		http.Error(w, "call_id is required", http.StatusBadRequest)
		return
	}

	if err := h.sipClient.HangupCall(req.CallID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "call ended"})
}

// CallStatus returns the status of a specific call
func (h *Handler) CallStatus(w http.ResponseWriter, r *http.Request) {
	callID := r.URL.Query().Get("call_id")
	if callID == "" {
		http.Error(w, "call_id parameter is required", http.StatusBadRequest)
		return
	}

	call, err := h.sipClient.GetCallStatus(callID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, call)
}

// ListExtensions returns available SIP extensions
func (h *Handler) ListExtensions(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would query FreeSWITCH for registered extensions
	extensions := []map[string]interface{}{
		{"extension": "1000", "status": "registered"},
		{"extension": "1001", "status": "registered"},
		{"extension": "1002", "status": "available"},
	}

	writeJSON(w, http.StatusOK, extensions)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
