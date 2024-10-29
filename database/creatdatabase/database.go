package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func parseDate(dateStr string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Erreur de parsing de la date:", err)
	}
	return t
}

func main() {
	data, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	// Create tables
	_, err = data.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts (
		post_id INTEGER PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		createdAt DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		content TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS likes (
		like_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		count INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS categories (
		category_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample data for users
	userStmt, err := data.Prepare(`INSERT INTO users (name, email, password) VALUES (?, ?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer userStmt.Close()

	users := []struct {
		name     string
		email    string
		password string
	}{
		{"Alice", "alice@example.com", "password123"},
		{"Bob", "bob@example.com", "password123"},
		{"Charlie", "charlie@example.com", "password123"},
		{"Diana", "diana@example.com", "password123"},
	}

	for _, user := range users {
		_, err = userStmt.Exec(user.name, user.email, user.password)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert sample data for posts
	postStmt, err := data.Prepare(`INSERT INTO posts (post_id, title, content, createdAt, user_id) VALUES (?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer postStmt.Close()

	posts := []struct {
		postID    int
		title     string
		content   string
		createdAt time.Time
		userID    int
	}{
		{1, "Premier Post", "Ceci est le contenu du premier post.", parseDate("2024-10-25 10:00:00"), 1},
		{2, "Deuxième Post", "Ceci est le contenu du deuxième post.", parseDate("2024-10-25 11:00:00"), 2},
		{3, "Troisième Post", "Ceci est le contenu du troisième post.", parseDate("2024-10-25 12:00:00"), 3},
	}

	for _, post := range posts {
		_, err = postStmt.Exec(post.postID, post.title, post.content, post.createdAt, post.userID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert sample data for comments
	commentStmt, err := data.Prepare(`INSERT INTO comments (post_id, content) VALUES (?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer commentStmt.Close()

	comments := []struct {
		postID  int
		content string
	}{
		{1, "Premier commentaire sur le premier post."},
		{1, "Deuxième commentaire sur le premier post."},
		{2, "Commentaire sur le deuxième post."},
	}

	for _, comment := range comments {
		_, err = commentStmt.Exec(comment.postID, comment.content)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert sample data for likes
	likeStmt, err := data.Prepare(`INSERT INTO likes (post_id, count) VALUES (?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer likeStmt.Close()

	likes := []struct {
		postID int
		count  int
	}{
		{1, 10},
		{2, 5},
		{3, 8},
	}

	for _, like := range likes {
		_, err = likeStmt.Exec(like.postID, like.count)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert sample data for categories
	categoryStmt, err := data.Prepare(`INSERT INTO categories (name) VALUES (?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer categoryStmt.Close()

	categories := []struct {
		name string
	}{
		{"Technology"},
		{"Health"},
		{"Education"},
		{"Travel"},
	}

	for _, category := range categories {
		_, err = categoryStmt.Exec(category.name)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Sample data inserted successfully!")
}
