package database

import (
	"chat/settings"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Migration() {
	log.Println("Старт миграции")
	db, err := settings.GetDatabaseConnect()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных %s\n", settings.PsqlDatabase)
	}
	defer db.Close()

	log.Printf("Успешное подключение к базе данных %s\n", settings.PsqlDatabase)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
		    id SERIAL PRIMARY KEY,
		    timestamp BIGINT,
			from_login VARCHAR(255) NOT NULL,
			to_login VARCHAR(255) NOT NULL,
			message TEXT NOT NULL
		)
`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
			login VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
`)

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Не удалось создать таблицы в базе данных %s\n", settings.PsqlDatabase)
	}
}
