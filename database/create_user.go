package database

import (
	"chat/settings"
	_ "github.com/lib/pq"
	"log"
)

func CreateUser(login, password string) error {
	db, err := settings.GetDatabaseConnect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)

	log.Printf("Создание пользователя %s\n", login)
	return err
}
