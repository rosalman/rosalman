
package internal

import (
	"encoding/json"
	"html/template"
	"net/http"
)



func DisplayPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the posts page for GET requests
		t, err := template.New("displayposts.html").ParseFiles("templates/displayposts.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch posts and their authors
		rows, err := db.Query(`
			SELECT p.post_id, p.title, p.content, u.username, p.created_at
			FROM posts p
			JOIN users u ON p.user_id = u.user_id
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

		// Render posts template
		err = t.Execute(w, map[string]interface{}{"Posts": posts})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle non-GET requests (POST or others)
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// If no GET method, continue with the original functionality:
	// Fetch posts and their authors
	rows, err := db.Query(`
        SELECT p.post_id, p.title, p.content, u.username, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.user_id
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
        SELECT c.comment_id, c.content, u.username, c.created_at
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
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
