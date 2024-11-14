package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a connection to the SQLite database
	data, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
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

CREATE TABLE IF NOT EXISTS categories (
	category_id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
	post_id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255) NOT NULL,
	content TEXT NOT NULL,
	createdAt DATETIME NOT NULL,
	user_id INTEGER,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts_categories (
	post_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
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
	user_id INTEGER,
	is_like INTEGER,
	type TEXT,
	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Insert test data
	_, err = data.Exec(`
-- Insert users
INSERT INTO users (name, email, password) VALUES 
    ('Alice', 'alice@example.com', 'password123'),
    ('Bob', 'bob@example.com', 'password456'),
    ('Charlie', 'charlie@example.com', 'password789');

-- Insert categories
INSERT INTO categories (name) VALUES 
    ('Technology'),
    ('Science'),
    ('Art');

-- Insert posts
INSERT INTO posts (title, content, createdAt, user_id) VALUES 
    ('First Post', 'This is the first test post.', datetime('now'), 1),
    ('Second Post', 'This is the second test post.', datetime('now'), 2),
    ('Third Post', 'This is the third test post.', datetime('now'), 3);

-- Link posts to categories
INSERT INTO posts_categories (post_id, category_id) VALUES 
    (1, 1),
    (2, 2),
    (3, 3);

-- Insert comments
INSERT INTO comments (post_id, user_id, content) VALUES 
    (1, 2, 'This is a comment on the first post by Bob.'),
    (2, 3, 'This is a comment on the second post by Charlie.'),
    (3, 1, 'This is a comment on the third post by Alice.');

-- Insert likes
INSERT INTO likes (post_id, user_id, is_like, type) VALUES 
    (1, 1, 1, 'post'),
    (1, 2, 1, 'post'),
    (2, 3, 0, 'post'),
    (3, 1, 1, 'post');
`)
	if err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}
}
