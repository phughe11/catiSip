package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Feedback represents user feedback
type Feedback struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

// FeedbackStore manages feedback in memory
type FeedbackStore struct {
	mu        sync.RWMutex
	feedbacks []Feedback
	nextID    int
}

var store = &FeedbackStore{
	feedbacks: []Feedback{},
	nextID:    1,
}

// enableCORS adds CORS headers to the response
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handleOptions handles preflight OPTIONS requests
func handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.WriteHeader(http.StatusOK)
}

// submitFeedback handles POST /api/feedback
func submitFeedback(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	var feedback Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if feedback.Name == "" || feedback.Email == "" || feedback.Message == "" {
		http.Error(w, "Name, email, and message are required", http.StatusBadRequest)
		return
	}

	// Validate rating
	if feedback.Rating < 1 || feedback.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	feedback.ID = fmt.Sprintf("feedback-%d", store.nextID)
	store.nextID++
	feedback.CreatedAt = time.Now()
	store.feedbacks = append(store.feedbacks, feedback)
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feedback)
}

// getFeedbacks handles GET /api/feedback
func getFeedbacks(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	store.mu.RLock()
	feedbacks := make([]Feedback, len(store.feedbacks))
	copy(feedbacks, store.feedbacks)
	store.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedbacks)
}

func main() {
	r := mux.NewRouter()

	// Feedback endpoints
	r.HandleFunc("/api/feedback", submitFeedback).Methods("POST")
	r.HandleFunc("/api/feedback", getFeedbacks).Methods("GET")
	r.HandleFunc("/api/feedback", handleOptions).Methods("OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
