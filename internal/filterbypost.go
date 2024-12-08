package internal

import (
	"encoding/json"
	"net/http"
	"text/template"
)

func FilterPostsByUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Use the correct path
        t, err := template.New("filter_posts.html").ParseFiles("templates/filterbypost.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
        return
    }

	// Handle API requests for filtered posts
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
		return
	}

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

	rows, err := db.Query(`
        SELECT id, title, content, created_at 
        FROM posts 
        WHERE user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var post struct {
			ID        int
			Title     string
			Content   string
			CreatedAt string
		}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Failed to parse post data", http.StatusInternalServerError)
			return
		}
		posts = append(posts, map[string]interface{}{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"created_at": post.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
