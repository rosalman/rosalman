package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("homepage.html").ParseFiles("templates/homepage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("login.html").ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var username, password string

	// Detect content type and parse accordingly
	if r.Header.Get("Content-Type") == "application/json" {
		var loginReq LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		username, password = loginReq.Username, loginReq.Password
	} else {
		// Default to form-encoded handling
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

	// Query the database for the user
	var user User // Use the User struct to hold the data
	query := `SELECT user_id, username, email, password, created_at FROM users WHERE username = ?`
	err = db.QueryRow(query, username).Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.CreationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Verify the password by comparing the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a new session token and set expiration time
	sessionToken := uuid.NewString()                 // Generate a new session token (UUID as string)
	expirationTime := time.Now().Add(24 * time.Hour) // Set session expiration to 24 hours

	// Insert session data into the sessions table
	_, err = db.Exec(`INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`,
		user.UserID, sessionToken, expirationTime)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true, // Ensures cookie is accessible only through HTTP requests (not JavaScript)
	})

	// Respond back to the client with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful")

	/*  sessionToken := uuid.NewString()
	expirationTime := time.Now().Add(24 * time.Hour)
	_, err = db.Exec(`INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`,
		userID, sessionToken, expirationTime)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful") */
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.New("signup.html").ParseFiles("templates/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the form values
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate the input
	if username == "" || email == "" || password == "" {
		http.Error(w, "Please fill out all fields", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Store the user data in the database
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err = db.Exec(query, username, email, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract category name from URL
	categoryName := strings.TrimPrefix(r.URL.Path, "/categories")

	// Build path to category template
	templatePath := fmt.Sprintf("templates/categories/%s.html", categoryName)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	t.Execute(w, nil)
}
