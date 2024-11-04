package database

import "fmt"

func CreateAcount(name, email, passeord string) error {
	_, err := db.Exec("INSERT INTO  users (name,email,password) VALUES ($1,$2,$3)", name, email, passeord)
	return err
}

func Createcomment(comment, post_id string) error {
	fmt.Print(comment)
	_, err := db.Exec("INSERT INTO comments (post_id,content) VALUES ($1,$2)", post_id, comment)
	return err
}
