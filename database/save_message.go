package database

import (
	"chat/settings"
	_ "github.com/lib/pq"
	"log"
)

func CreateMessage(timestamp int64, from, to, message string) error {
	db, err := settings.GetDatabaseConnect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO messages (timestamp, from_login, to_login, message) VALUES ($1, $2, $3, $4)", timestamp, from, to, message)

	log.Printf("%d Новое сообщение от %s для %s: %s", timestamp, from, to, message)
	return err
}
