package main

import (
	"log"
	"net/http"
	"os"

	"github.com/phughe11/catiSip/backend/config"
	"github.com/phughe11/catiSip/backend/handlers"
	"github.com/phughe11/catiSip/backend/sip"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize SIP client
	sipClient, err := sip.NewClient(cfg.SIP)
	if err != nil {
		log.Fatalf("Failed to initialize SIP client: %v", err)
	}
	defer sipClient.Close()

	// Setup HTTP handlers
	mux := http.NewServeMux()
	h := handlers.New(sipClient, cfg)

	// API routes
	mux.HandleFunc("/api/health", h.HealthCheck)
	mux.HandleFunc("/api/call/make", h.MakeCall)
	mux.HandleFunc("/api/call/hangup", h.HangupCall)
	mux.HandleFunc("/api/call/status", h.CallStatus)
	mux.HandleFunc("/api/extensions", h.ListExtensions)

	// CORS middleware
	corsHandler := enableCORS(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
