package settings

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getDataFromEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке переменных окружения: %v", err)
	}

	return os.Getenv(key)
}

var PsqlHost = getDataFromEnv("PSQL_HOST")
var PsqlPort = getDataFromEnv("PSQL_PORT")
var PsqlUser = getDataFromEnv("PSQL_USER")
var PsqlPassword = getDataFromEnv("PSQL_PASSWORD")
var PsqlDatabase = getDataFromEnv("PSQL_DATABASE")
var ChatHost = getDataFromEnv("CHAT_HOST")
var ChatPort = getDataFromEnv("CHAT_PORT")

func GetDatabaseConnect() (*sql.DB, error) {
	connStr := " host=" + PsqlHost +
		" port=" + PsqlPort +
		" user=" + PsqlUser +
		" password=" + PsqlPassword +
		" dbname=" + PsqlDatabase +
		" sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	return db, err
}
