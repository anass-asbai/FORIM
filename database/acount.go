package database

func CreateAcount(name, email, passeord string) error {
	_, err := db.Exec("INSERT INTO  users (name,email,password) VALUES ($1,$2,$3)", name, email, passeord)
	return err
}
