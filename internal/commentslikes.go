package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func ReactToCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// Render the postlikes.html page for GET requests
		t, err := template.ParseFiles("templates/postlikes.html")
		if err != nil {
			http.Error(w, "Failed to load the HTML page", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check user session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session token
	var userID int
	err = db.QueryRow(`
        SELECT user_id 
        FROM sessions 
        WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP`,
		sessionCookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
		return
	}

	// Define a struct for the reaction
	var reaction struct {
		CommentID int    `json:"comment_id"`
		Reaction  string `json:"reaction"` // "like" or "dislike"
	}

	// Handle Content-Type: application/json
	if r.Header.Get("Content-Type") == "application/json" {
		// Parse JSON request body
		err := json.NewDecoder(r.Body).Decode(&reaction)
		if err != nil || (reaction.Reaction != "like" && reaction.Reaction != "dislike") {
			http.Error(w, "Invalid request payload or reaction type", http.StatusBadRequest)
			return
		}
	} else {
		// Handle form-encoded data (application/x-www-form-urlencoded)
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Parse form values
		reaction.CommentID, err = strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		reaction.Reaction = r.FormValue("reaction")
		if reaction.Reaction != "like" && reaction.Reaction != "dislike" {
			http.Error(w, "Invalid reaction type", http.StatusBadRequest)
			return
		}
	}

	// Insert or update the reaction
	_, err = db.Exec(`
        INSERT INTO comment_reactions (comment_id, user_id, reaction)
        VALUES (?, ?, ?)
        ON CONFLICT(comment_id, user_id) DO UPDATE SET reaction = excluded.reaction`,
		reaction.CommentID, userID, reaction.Reaction)
	if err != nil {
		http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reaction recorded successfully"})
}
