package database

import (
	pb "chat/gen/pb-go/chat/protos"
	"chat/settings"
	_ "github.com/lib/pq"
	"log"
)

func GetMessages(from, to string) ([]*pb.ChatMessage, error) {
	db, err := settings.GetDatabaseConnect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT timestamp, from_login, to_login, message FROM messages WHERE (from_login = $1 AND to_login = $2) OR (from_login = $2 AND to_login = $1) ORDER BY timestamp", from, to)
	if err != nil {
		log.Printf("Не удалось получить сообщения от %s для %s", from, to)
		return nil, err
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var timestamp int64
		var fromLogin, toLogin, message string
		err = rows.Scan(&timestamp, &fromLogin, &toLogin, &message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &pb.ChatMessage{
			Timestamp: timestamp,
			From:      fromLogin,
			To:        toLogin,
			Message:   message,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Не удалось получить сообщения от %s для %s", from, to)
		return nil, err
	}

	log.Printf("Успешное получение сообщений от %s для %s", from, to)
	return messages, nil
}
