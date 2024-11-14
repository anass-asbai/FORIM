package database

import (
	"fmt"
	"strconv"
	"time"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Date     time.Time
	Like     int
	Deslike  int
	User     string
	Category string
}
type Comment struct {
	ID      int
	PostID  int
	Comment string
	User    string
}

func InsertPost(title, content, email string, categories []string) error {
	user_id := 0

	erre := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&user_id)
	if erre != nil {
		return erre
	}

	result, err := db.Exec("INSERT INTO posts (title, content, createdAt, user_id) VALUES (?, ?, datetime('now'), ?)", title, content, user_id)
	if err != nil {
		return err
	}

	idPost, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, category := range categories {
		category_id := 0
		err := db.QueryRow("SELECT category_id FROM categories WHERE name = ?", category).Scan(&category_id)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)", idPost, category_id)
		if err != nil {
			return err
		}
	}

	return nil
}

func CountPost(limit int) bool {
	var l string
	err := db.QueryRow("SELECT COUNT(post_id) FROM posts").Scan(&l)
	if err != nil {
		fmt.Println("Error querying row:", err)
		return false
	}

	max, err := strconv.Atoi(l)
	if err != nil {
		fmt.Println("Error converting count to integer:", err)
		return false
	}

	return max > limit
}

func GetPosts(catigorie string, limit int) ([]Post, error) {

	if limit < 0 {
		limit = 0
	}
	var query string
	query = `SELECT 
    posts.post_id,
    posts.title,
    posts.content,
    posts.createdAt,
    GROUP_CONCAT(COALESCE(categories.name,'')) AS category_name,
    COALESCE(users.name, '') AS username,
    COALESCE(SUM(CASE WHEN likes.is_like = 1 THEN 1 ELSE 0 END), 0) AS like_count,
    COALESCE(SUM(CASE WHEN likes.is_like = 2 THEN 1 ELSE 0 END), 0) AS dislike_count
FROM 
    posts
LEFT JOIN users ON users.user_id = posts.user_id
LEFT JOIN likes ON likes.post_id = posts.post_id
LEFT JOIN posts_categories ON posts_categories.post_id = posts.post_id
LEFT JOIN categories ON categories.category_id = posts_categories.category_id
GROUP BY 
    posts.post_id
`

	if catigorie != "" {
		query += " WHERE category = ?"
	}
	query += " ORDER BY posts.createdAt DESC LIMIT 5 OFFSET " + strconv.Itoa(limit)
	rows, err := db.Query(query, catigorie)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.Category, &post.User, &post.Like, &post.Deslike); err != nil {
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

func InsertLike(id, email string, is_like bool) error {
	var id_user int
	var checkrow int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id_user)
	if err != nil {
		return err
	}
	pre, err := db.Prepare("SELECT COUNT(like_id) FROM likes WHERE post_id = ? AND user_id = ? ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer pre.Close()
	err = pre.QueryRow(id, id_user).Scan(&checkrow)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if checkrow == 0 && is_like {
		if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 1); err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else if checkrow == 0 && !is_like {
		if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 2); err != nil {
			fmt.Println(err.Error())
			return err
		}

	} else if checkrow != 0 && is_like {
		var reaction int
		pre, err := db.Prepare("SELECT is_like FROM likes WHERE post_id = ? AND user_id = ? ")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer pre.Close()
		err = pre.QueryRow(id, id_user).Scan(&reaction)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if _, err := db.Exec(`DELETE FROM Likes WHERE post_id = ? AND user_id = ? `, id, id_user); err != nil {
			fmt.Println(err.Error())
			return err
		}
		if reaction == 2 {
			if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 1); err != nil {
				fmt.Println(err.Error())
				return err
			}
		}

	} else if checkrow != 0 && !is_like {

		var reaction int
		pre, err := db.Prepare("SELECT is_like FROM likes WHERE post_id = ? AND user_id = ? ")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer pre.Close()
		err = pre.QueryRow(id, id_user).Scan(&reaction)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if _, err := db.Exec(`DELETE FROM Likes WHERE post_id = ? AND user_id = ? `, id, id_user); err != nil {
			fmt.Println(err.Error())
			return err
		}
		if reaction == 1 {
			if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 2); err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}
	return nil
}
