package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phughe11/catiSip/backend/config"
	"github.com/phughe11/catiSip/backend/sip"
)

func setupTest() *Handler {
	cfg := &config.Config{
		SIP: config.SIPConfig{
			Host:     "localhost",
			Port:     5060,
			Username: "1000",
			Password: "1234",
			Domain:   "localhost",
		},
		Server: config.ServerConfig{
			Port: 8080,
		},
	}

	sipClient, _ := sip.NewClient(cfg.SIP)
	return New(sipClient, cfg)
}

func TestHealthCheck(t *testing.T) {
	h := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestMakeCall(t *testing.T) {
	h := setupTest()

	reqBody := map[string]string{
		"from": "1000",
		"to":   "1001",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/call/make", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.MakeCall(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response sip.Call
	json.NewDecoder(w.Body).Decode(&response)

	if response.From != "1000" || response.To != "1001" {
		t.Errorf("Expected from=1000 to=1001, got from=%s to=%s", response.From, response.To)
	}

	if response.Status != "dialing" {
		t.Errorf("Expected status 'dialing', got '%s'", response.Status)
	}
}

func TestMakeCallInvalidMethod(t *testing.T) {
	h := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/call/make", nil)
	w := httptest.NewRecorder()

	h.MakeCall(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestMakeCallMissingFields(t *testing.T) {
	h := setupTest()

	reqBody := map[string]string{
		"from": "1000",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/call/make", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.MakeCall(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListExtensions(t *testing.T) {
	h := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/api/extensions", nil)
	w := httptest.NewRecorder()

	h.ListExtensions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) == 0 {
		t.Error("Expected at least one extension, got none")
	}
}
