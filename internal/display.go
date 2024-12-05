package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func DisplayPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
            SELECT p.id, p.title, p.content, u.username, p.created_at
            FROM posts p
            JOIN users u ON p.author_id = u.id
            ORDER BY p.created_at DESC
        `)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		posts := []PostResponse{}
		for rows.Next() {
			var post PostResponse
			rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)

			// Fetch comments for the post
			commentRows, _ := db.Query(`
                SELECT c.id, c.text, u.username, c.created_at
                FROM comments c
                JOIN users u ON c.author_id = u.id
                WHERE c.post_id = ?
            `, post.ID)

			comments := []CommentResponse{}
			for commentRows.Next() {
				var comment CommentResponse
				commentRows.Scan(&comment.ID, &comment.Text, &comment.Author, &comment.CreatedAt)
				comments = append(comments, comment)
			}
			commentRows.Close()

			post.Comments = comments
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
