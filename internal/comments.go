package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check authentication
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

		var commentReq CommentRequest

		// Detect content type and parse accordingly
		if r.Header.Get("Content-Type") == "application/json" {
			err = json.NewDecoder(r.Body).Decode(&commentReq)
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
			commentReq.PostID, _ = strconv.Atoi(r.FormValue("post_id"))
			commentReq.Comment = r.FormValue("comment")
		}

		if commentReq.PostID == 0 || commentReq.Comment == "" {
			http.Error(w, "Post ID and comment text are required", http.StatusBadRequest)
			return
		}

		// Insert the comment into the database
		_, err = db.Exec(`INSERT INTO comments (text, post_id, author_id) VALUES (?, ?, ?)`,
			commentReq.Comment, commentReq.PostID, userID)
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
	}
}