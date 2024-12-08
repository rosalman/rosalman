package internal

import "time"

type User struct {
    UserID       int64     `json:"user_id"`         // Adding the UserID field
    Username     string    `json:"username"`        // Username with JSON tag
    Email        string    `json:"email"`           // Email with JSON tag
    PasswordHash string    `json:"password_hash"`   // Password hash with a corrected name and JSON tag
    CreationDate time.Time `json:"creation_date"`   // CreationDate with JSON tag
    Password     string    `json:"password"`        // Password field for form-based requests
    SecurityAnswer string  `json:"securityQuestion"` // Security question answer
}


type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
    SessionID string    `json:"session_id"`  // SessionID will be a unique string (UUID)
    UserID    int       `json:"user_id"`     // UserID corresponds to the user in the session
    Username  string    `json:"username"`    // Username of the logged-in user
    Expires   time.Time `json:"expires"`     // Expiration time of the session
}
type PostRequest struct {
    Title      string   `json:"title"`
    Content    string   `json:"content"`
    Categories []int    `json:"categories"` // IDs of associated categories
}

type CommentRequest struct {
    PostID  int    `json:"post_id"`
    Comment string `json:"comment"`
}

type PostResponse struct {
	ID        int               `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Author    string            `json:"author"`
	CreatedAt string            `json:"created_at"`
	Likes     int               `json:"likes"`
	Dislikes  int               `json:"dislikes"`
	Comments  []CommentResponse `json:"comments"`
}

type CommentResponse struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
}