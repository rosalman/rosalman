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

	// Serve static files for CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
