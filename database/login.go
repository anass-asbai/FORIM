package database

import "forim/bcryptp"

func GetLogin(email, password string) (bool, error) {
	rows, err := db.Query("SELECT password FROM users WHERE email =  $1", email)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {

		var passwords string
		if err := rows.Scan(&passwords); err != nil {
			return false, err
		}
		if bcryptp.CheckPasswordHash(password, passwords) {
			return true, nil
		}

	}
	return false, nil
}
