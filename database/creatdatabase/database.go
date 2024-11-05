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
    user_id INTEGER, 
    content TEXT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE  
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

	// Insert sample data
	_, err = data.Exec(`
		INSERT INTO users (name, email, password) VALUES
		('Alice', 'alice@example.com', 'password123'),
		('Bob', 'bob@example.com', 'password123'),
		('Charlie', 'charlie@example.com', 'password123');
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = data.Exec(`
		INSERT INTO posts (title, content, createdAt, user_id) VALUES
		('First Post', 'This is the content of the first post.', '2024-11-01 12:00:00', 1),
		('Second Post', 'This is the content of the second post.', '2024-11-02 12:00:00', 2),
		('Third Post', 'This is the content of the third post.', '2024-11-03 12:00:00', 3);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = data.Exec(`
		INSERT INTO comments (post_id, user_id, content) VALUES
		(1, 2, 'Great post!'),
		(1, 3, 'Thanks for sharing.'),
		(2, 1, 'Interesting perspective.');
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = data.Exec(`
		INSERT INTO likes (post_id, count) VALUES
		(1, 5),
		(2, 3),
		(3, 10);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = data.Exec(`
		INSERT INTO categories (name) VALUES
		('Technology'),
		('Lifestyle'),
		('Education');
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Données test insérées avec succès.")

}
