package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check for authentication
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
			return
		}

		var userID int
		err = db.QueryRow(`SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP`,
			sessionCookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}

		var postReq PostRequest

		// Detect content type and parse accordingly
		if r.Header.Get("Content-Type") == "application/json" {
			err = json.NewDecoder(r.Body).Decode(&postReq)
			if err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}
		} else {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}
			postReq.Title = r.FormValue("title")
			postReq.Content = r.FormValue("content")
			// Categories can be passed as comma-separated values
		}

		if postReq.Title == "" || postReq.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Insert post into the database
		result, err := db.Exec(`INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)`,
			postReq.Title, postReq.Content, userID)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		postID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
			return
		}

		// Associate categories with the post
		for _, categoryID := range postReq.Categories {
			_, err = db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
				postID, categoryID)
			if err != nil {
				http.Error(w, "Failed to associate categories", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
	}
}