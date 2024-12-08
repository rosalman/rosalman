package internal

import (
	"encoding/json"
	"net/http"
	"text/template"
)

func FilterLikedPostsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// Serve the HTML page for filtering liked posts
		t, err := template.New("filterbylikes.html").ParseFiles("templates/filterbylikes.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

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
        SELECT p.id, p.title, p.content, p.created_at 
        FROM posts p
        INNER JOIN post_reactions pr ON p.id = pr.post_id
        WHERE pr.user_id = ? AND pr.reaction = 'like'`, userID)
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
