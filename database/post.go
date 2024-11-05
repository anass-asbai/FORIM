package database

import (
	"fmt"
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
	Like    int
	User    string
}
type Comment struct {
	ID      int
	PostID  int
	Comment string
	User    string
}

func InsertPost(title, content, email string) error {
	id := 0
	erre := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id)
	_ = erre
	fmt.Println(id)
	_, err := db.Exec("INSERT INTO posts (title, content, createdAt,user_id) VALUES (?, ?, datetime('now'),?)", title, content, id)
	return err
}

/*
func InsertComment(comment,email string) error{
	id := 0
	erre := db.QueryRow("SELECT FROM users WHERE email = ?",email).Scan(&id)
	return err
}*/

func GetPosts() ([]Post, error) {
	rows, err := db.Query(`SELECT posts.post_id, posts.title, posts.content, posts.createdAt, 
       COALESCE(likes.count, 0) AS count,COALESCE(users.name,'') AS username
       FROM posts 
       LEFT JOIN likes ON likes.post_id = posts.post_id 
	   LEFT JOIN users ON users.user_id=posts.user_id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.Like, &post.User); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetComment(id string) ([]Comment, error) {
	rows, err := db.Query(`SELECT content, users.name FROM comments INNER JOIN users ON users.user_id=comments.name WHERE  post_id = ?`, id)
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
		fmt.Println(rows.Next())
		Comments = append(Comments, Comment)
	}
	return Comments, nil
}
