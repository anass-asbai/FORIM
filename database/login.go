package database

func GetLogin(email, password string) (bool, error) {
	rows, err := db.Query("SELECT email,password FROM users")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var emails string
		var passwords string
		if err := rows.Scan(&emails, &passwords); err != nil {
			return false, err
		}
		if emails == email && passwords == password {
			return true, nil
		}

	}
	return false, nil
}
