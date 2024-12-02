package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}

	// Queries to create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id)
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS postCategories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (category_id) REFERENCES categories(category_id)
		)`,
		`CREATE TABLE IF NOT EXISTS likes (
			like_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			comment_id INTEGER,
			is_like BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (comment_id) REFERENCES comments(comment_id)
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			session_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			session_token TEXT NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		)`,
	}

	// Execute each query
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			fmt.Println("Error executing query:", query, "\nError:", err)
			return
		}
	}
}