package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	sessionTTL time.Duration
	tpl        *template.Template
	uploadDir  = "./uploads"
)

// InitHandlers now accepts the sessionTTL
func InitHandlers(ttl time.Duration) {
	sessionTTL = ttl

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Parse all templates
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// A new type for our context keys to avoid collisions
type contextKey string

const userKey contextKey = "userID"
const usernameKey contextKey = "username"

// authMiddleware protects routes that require a logged-in user
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No cookie, redirect to login
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		sessionToken := cookie.Value

		// 2. Get the session from the database
		session, err := GetSession(sessionToken)
		if err != nil {
			// Invalid session (not found, expired), delete cookie and redirect
			http.SetCookie(w, &http.Cookie{
				Name: "session_token", Value: "", Expires: time.Unix(0, 0), Path: "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// 3. Session is valid. Store user info in context for next handlers.
		ctx := context.WithValue(r.Context(), userKey, session.UserID)
		ctx = context.WithValue(ctx, usernameKey, session.Username)

		// 4. Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleLoginPage shows the login page or redirects to dashboard if already logged in
func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user already has a valid session cookie
	if cookie, err := r.Cookie("session_token"); err == nil {
		if _, err := GetSession(cookie.Value); err == nil {
			// Valid session exists, redirect to dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	// No valid session, show login page
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// handleLogin processes the login form
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := GetUserByUsername(username)
	if err != nil || !CheckPasswordHash(password, user.PasswordHash) {
		// Use HX-Redirect for HTMX-driven form submissions
		w.Header().Set("HX-Redirect", "/error")
		w.WriteHeader(http.StatusOK) // HTMX expects 200 on form submit
		return
	}

	// 1. Create a new session in the database
	newSession, err := CreateSession(user.ID, user.Username, sessionTTL)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		w.Header().Set("HX-Redirect", "/error")
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. Set the session token in an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    newSession.SessionToken,
		Expires:  newSession.ExpiresAt,
		Path:     "/",  // Available to entire site
		HttpOnly: true, // Not accessible via JavaScript
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// handleLogout clears the session
func handleLogout(w http.ResponseWriter, r *http.Request) {
	// 1. Get token from cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// 2. Delete session from database
		DeleteSession(cookie.Value)
	}

	// 3. Expire the cookie in the browser
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire now
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleErrorPage shows a simple login error
func handleErrorPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "error.html", nil)
}

// handleDashboard shows the main todo list
func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Get user info from context (set by authMiddleware)
	userID := r.Context().Value(userKey).(primitive.ObjectID)
	username := r.Context().Value(usernameKey).(string)

	todos, err := GetTodosByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": username,
		"Todos":    todos,
	}
	tpl.ExecuteTemplate(w, "dashboard.html", data)
}

// handleCreateTodo handles the creation of a new todo
func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user info from context
	userID := r.Context().Value(userKey).(primitive.ObjectID)

	// Max 10MB upload
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	// --- File Upload Logic ---
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename
	ext := filepath.Ext(handler.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file data
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	// --- End File Upload Logic ---

	todo := Todo{
		Title:     title,
		Content:   content,
		ImagePath: fileName, // Save only the filename
		UserID:    userID,
	}

	newTodo, err := CreateTodo(todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// Return the new todo item partial
	tpl.ExecuteTemplate(w, "_todo-item.html", newTodo)
}

// handleGetEditForm returns the HTML form for editing a todo
func handleGetEditForm(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value(userKey).(primitive.ObjectID)
	id := r.PathValue("id")

	todo, err := GetTodoByID(id, userID)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	tpl.ExecuteTemplate(w, "_edit-form.html", todo)
}

// handleUpdateTodo updates a todo's data
func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value(userKey).(primitive.ObjectID)
	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")

	if err := UpdateTodo(id, userID, title, content); err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	// Get the updated todo to send back
	updatedTodo, err := GetTodoByID(id, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated todo", http.StatusInternalServerError)
		return
	}

	// Return the updated todo item partial
	tpl.ExecuteTemplate(w, "_todo-item.html", updatedTodo)
}

// handleGetTodoItem returns the standard view partial for a todo (used for cancel)
func handleGetTodoItem(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value(userKey).(primitive.ObjectID)
	id := r.PathValue("id")

	todo, err := GetTodoByID(id, userID)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	tpl.ExecuteTemplate(w, "_todo-item.html", todo)
}

// handleDeleteTodo deletes a todo
func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get user info from context
	userID := r.Context().Value(userKey).(primitive.ObjectID)
	id := r.PathValue("id")

	// Before deleting from DB, get the todo to find its image path
	todo, err := GetTodoByID(id, userID)
	if err != nil {
		log.Printf("Error finding todo for deletion: %v", err)
		// Continue to delete DB entry anyway
	}

	// Delete from database
	if err := DeleteTodo(id, userID); err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	// Delete the associated image file
	if todo != nil && todo.ImagePath != "" {
		filePath := filepath.Join(uploadDir, todo.ImagePath)
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to delete image file %s: %v", filePath, err)
		}
	}

	// Return an empty response, HTMX will remove the element
	w.WriteHeader(http.StatusOK)
}
