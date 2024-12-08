package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Retrieve the PostID (this should ideally be dynamic from URL params)
		postIDStr := r.URL.Query().Get("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID == 0 {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		// Render the comment creation page for GET requests and pass the PostID to the template
		t, err := template.New("create_comment.html").ParseFiles("templates/create_comment.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Pass the PostID to the template
		t.Execute(w, map[string]interface{}{
			"PostID": postID, // pass the actual post ID when rendering the form
		})
		return
	}

	if r.Method != http.MethodPost {
		// Reject non-POST requests with an error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var commentReq CommentRequest
	var userID int

	// Check authentication
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
		return
	}

	// Validate session
	query := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP`
	err = db.QueryRow(query, sessionCookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
		return
	}

	// Detect content type and parse accordingly
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&commentReq)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
	} else {
		// Default to form-encoded handling
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		// Extract the PostID and comment text from form values
		postID, _ := strconv.Atoi(r.FormValue("post_id"))
		commentReq = CommentRequest{
			PostID:  postID,
			Comment: r.FormValue("comment"),
		}
	}

	// Validate input
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

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}
