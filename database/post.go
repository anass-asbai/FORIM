package database

import (
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
	Like    int
	Comment string
	User    string
}

func InsertPost(title, content string) error {
	_, err := db.Exec("INSERT INTO posts (title, content, createdAt) VALUES (?, ?, datetime('now'))", title, content)
	return err
}

func GetPosts() ([]Post, error) {
	rows, err := db.Query(`SELECT posts.post_id, posts.title, posts.content, posts.createdAt, 
       COALESCE(likes.count, 0) AS count,  COALESCE(comments.content, '') AS comments,COALESCE(users.name,'') AS username
       FROM posts 
       LEFT JOIN likes ON likes.post_id = posts.post_id 
       LEFT JOIN comments ON comments.post_id = posts.post_id
	   LEFT JOIN users ON users.user_id=posts.user_id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.Like, &post.Comment, &post.User); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
