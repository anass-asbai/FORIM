package database

import (
	"fmt"
	"time"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Date     time.Time
	Like     int
	User     string
	Category string
}
type Comment struct {
	ID      int
	PostID  int
	Comment string
	User    string
}

func InsertPost(title, content, email, categories string) error {
	id := 0
	category_id := 0
	erre := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id)
	_ = erre
	erre = db.QueryRow("SELECT category_id FROM categories WHERE name = ?", categories).Scan(&category_id)
	fmt.Println(categories)
	_ = erre
	_, err := db.Exec("INSERT INTO posts (title, content, createdAt,user_id, category_id) VALUES (?, ?, datetime('now'),?,?)", title, content, id, category_id)
	return err
}

/*
func InsertComment(comment,email string) error{
	id := 0
	erre := db.QueryRow("SELECT FROM users WHERE email = ?",email).Scan(&id)
	return err
}*/

func GetPosts(catigorie string) ([]Post, error) {
	var query string
	query = `SELECT posts.post_id, posts.title, posts.content, posts.createdAt, 
         COALESCE(likes.count, 0) AS count, COALESCE(users.name, '') AS username,
         categories.name AS category
         FROM posts 
         LEFT JOIN likes ON likes.post_id = posts.post_id 
         LEFT JOIN users ON users.user_id = posts.user_id
         LEFT JOIN categories ON categories.category_id = posts.category_id`

	if catigorie != "" {
		query += " WHERE category = ?"
	}
	query += " ORDER BY likes.count DESC;"
	rows, err := db.Query(query, catigorie)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.Like, &post.User, &post.Category); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetComment(id string) ([]Comment, error) {
	rows, err := db.Query(`SELECT content, users.name FROM comments LEFT Join users ON users.user_id=comments.user_id WHERE  post_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Comments []Comment
	for rows.Next() {

		var Comment Comment
		if err := rows.Scan(&Comment.Comment, &Comment.User); err != nil {
			return nil, err
		}

		Comments = append(Comments, Comment)
	}
	return Comments, nil
}
