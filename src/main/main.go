package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default
	}

	// Get Session TTL from env, default to 10 minutes
	ttlMinutes, err := strconv.Atoi(os.Getenv("SESSION_TTL_MINUTES"))
	if err != nil || ttlMinutes <= 0 {
		ttlMinutes = 10
	}
	sessionTTL := time.Duration(ttlMinutes) * time.Minute

	// Initialize Database
	InitDB(mongoURI)
	CreateAdminUser() // Create admin user on startup

	// Initialize Handlers, passing in the TTL
	InitHandlers(sessionTTL)

	// --- Setup Router ---
	mux := http.NewServeMux()

	// --- Public Routes ---

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		version, _ := os.ReadFile("version")
		response := map[string]string{"status": "ok", "version": string(version)}
		json.NewEncoder(w).Encode(response)
	})

	// Static file server for images
	// We specify "GET" to make the pattern non-conflicting (Go 1.22+)
	fs := http.FileServer(http.Dir(uploadDir))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fs))

	// Public page handlers
	mux.HandleFunc("GET /", handleLoginPage)
	mux.HandleFunc("POST /login", handleLogin)
	mux.HandleFunc("GET /logout", handleLogout)
	mux.HandleFunc("GET /error", handleErrorPage)

	// --- Protected Routes ---
	// We apply the authMiddleware wrapper to each
	// handler function individually.
	mux.Handle("GET /dashboard", authMiddleware(http.HandlerFunc(handleDashboard)))
	mux.Handle("POST /todos", authMiddleware(http.HandlerFunc(handleCreateTodo)))
	mux.Handle("GET /todos/{id}/edit", authMiddleware(http.HandlerFunc(handleGetEditForm)))
	mux.Handle("PUT /todos/{id}", authMiddleware(http.HandlerFunc(handleUpdateTodo)))
	mux.Handle("GET /todos/{id}/item", authMiddleware(http.HandlerFunc(handleGetTodoItem)))
	mux.Handle("DELETE /todos/{id}", authMiddleware(http.HandlerFunc(handleDeleteTodo)))

	// Start server
	log.Printf("Server starting on :8080 (Session TTL: %s)", sessionTTL)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
