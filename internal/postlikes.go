package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
)

func ReactToPostHandler(w http.ResponseWriter, r *http.Request) {
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

	var userID int
	err = db.QueryRow(`SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP`,
		sessionCookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
		return
	}

	// Parse form-encoded data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	reaction := r.FormValue("reaction")
	if reaction != "like" && reaction != "dislike" {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	// Insert or update the reaction in the database
	_, err = db.Exec(`
        INSERT INTO post_reactions (post_id, user_id, reaction)
        VALUES (?, ?, ?)
        ON CONFLICT(post_id, user_id) DO UPDATE SET reaction = excluded.reaction`,
		postID, userID, reaction)
	if err != nil {
		http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reaction recorded successfully"})
}

/* func ReactToPostHandler (w http.ResponseWriter, r *http.Request) {
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

		var userID int
		err = db.QueryRow(`SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP`,
			sessionCookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var reaction struct {
			PostID  int    `json:"post_id"`
			Reaction string `json:"reaction"` // "like" or "dislike"
		}
		err = json.NewDecoder(r.Body).Decode(&reaction)
		if err != nil || (reaction.Reaction != "like" && reaction.Reaction != "dislike") {
			http.Error(w, "Invalid request payload or reaction type", http.StatusBadRequest)
			return
		}

		// Insert or update the reaction
		_, err = db.Exec(`
			INSERT INTO post_reactions (post_id, user_id, reaction)
			VALUES (?, ?, ?)
			ON CONFLICT(post_id, user_id) DO UPDATE SET reaction = excluded.reaction`,
			reaction.PostID, userID, reaction.Reaction)
		if err != nil {
			http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Reaction recorded successfully"})
	}
} */

func FetchPostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// Use PostResponse instead of defining a new struct
	var post PostResponse

	err = db.QueryRow(`
			SELECT p.id, p.title, p.content, u.username, p.created_at,
				(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction = 'like') AS likes,
				(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction = 'dislike') AS dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.user_id
			WHERE p.id = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	// Fetch comments
	rows, err := db.Query(`
SELECT c.comment_id, c.content, u.username, c.created_at,
	(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.comment_id AND reaction = 'like') AS likes,
	(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.comment_id AND reaction = 'dislike') AS dislikes
FROM comments c
JOIN users u ON c.user_id = u.user_id
WHERE c.post_id = ?
ORDER BY c.created_at ASC`, postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment CommentResponse
		err := rows.Scan(&comment.ID, &comment.Text, &comment.Author, &comment.CreatedAt, &comment.Likes, &comment.Dislikes)
		if err != nil {
			http.Error(w, "Failed to parse comment", http.StatusInternalServerError)
			return
		}
		post.Comments = append(post.Comments, comment)
	}
}
