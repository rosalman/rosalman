package internal

import (
	"encoding/json"
	"net/http"
)

func DisplayPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Fetch posts and their authors
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Error processing post data", http.StatusInternalServerError)
			return
		}

		// Fetch comments for each post
		comments, err := fetchComments(post.ID)
		if err != nil {
			http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			return
		}

		post.Comments = comments
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Helper function to fetch comments for a given post ID
func fetchComments(postID int) ([]CommentResponse, error) {
	commentRows, err := db.Query(`
        SELECT c.id, c.text, u.username, c.created_at
        FROM comments c
        JOIN users u ON c.author_id = u.id
        WHERE c.post_id = ?
    `, postID)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	comments := []CommentResponse{}
	for commentRows.Next() {
		var comment CommentResponse
		err := commentRows.Scan(&comment.ID, &comment.Text, &comment.Author, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

/* func DisplayPostsHandler(db *sql.DB) http.HandlerFunc {
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
} */
