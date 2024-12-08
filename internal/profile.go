package internal

import (
	"net/http"
	"text/template"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
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

	// Fetch the user details using the userID
	var user User
	err = db.QueryRow(`
        SELECT user_id, username, email, created_at 
        FROM users 
        WHERE user_id = ?`, userID).Scan(&user.UserID, &user.Username, &user.Email, &user.CreationDate)
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	// Serve the profile page
	if r.Method == http.MethodGet {
		// Serve the HTML page for the user's profile
		t, err := template.New("profile.html").ParseFiles("templates/profile.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Pass user details to the template
		err = t.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle other methods (if necessary)
	// You could allow profile updates (POST) here if you choose.
}
