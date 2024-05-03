package database

import (
	"chat/settings"
	_ "github.com/lib/pq"
	"log"
)

func AuthUser(login, password string) (bool, error) {
	db, err := settings.GetDatabaseConnect()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1 AND password = $2", login, password).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		log.Printf("Успешная аутентификация пользователя %s\n", login)
		return true, nil
	}

	log.Printf("Неудачная аутентификация пользователя %s\n", login)
	return false, nil
}
