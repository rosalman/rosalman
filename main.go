package main

import (
	"fmt"
	internal "forum/internal"
	"log"
	"net/http"
)

func main() {
	internal.Init()

	// routes
	http.HandleFunc("/", internal.HomepageHandler)
	http.HandleFunc("/login", internal.LoginHandler)
	http.HandleFunc("/signup", internal.SignupHandler)
	http.HandleFunc("/animals", internal.CategoryHandler)
	http.HandleFunc("/art", internal.CategoryHandler)
	http.HandleFunc("/beauty", internal.CategoryHandler)
	http.HandleFunc("/finance", internal.CategoryHandler)
	http.HandleFunc("/fitness", internal.CategoryHandler)
	http.HandleFunc("/food", internal.CategoryHandler)
	http.HandleFunc("/health", internal.CategoryHandler)
	http.HandleFunc("/jobs", internal.CategoryHandler)
	http.HandleFunc("/music", internal.CategoryHandler)
	http.HandleFunc("/sport", internal.CategoryHandler)
	http.HandleFunc("/technology", internal.CategoryHandler)
	http.HandleFunc("/travel", internal.CategoryHandler)

	// New Routes for Posts and Comments
	http.HandleFunc("/posts/create", internal.CreatePostHandler)
	http.HandleFunc("/comments/create", internal.CreateCommentHandler)
	http.HandleFunc("/posts", internal.DisplayPostsHandler)

	http.HandleFunc("/react-to-post", internal.ReactToPostHandler)
	http.HandleFunc("/react-to-comment", internal.ReactToCommentHandler)
	http.HandleFunc("/fetch-post-details", internal.FetchPostDetailsHandler)
	
	  // Route for filtering liked posts
	http.HandleFunc("/filter-by-liked-posts", internal.FilterLikedPostsHandler)
	   // Route for filtering posts by the logged-in user
	http.HandleFunc("/filter-by-user-posts", internal.FilterPostsByUserHandler)

	http.HandleFunc("/profile", internal.ProfileHandler)


	// Serve static files for CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
