package main

import (
	"bufio"
	pb "chat/gen/pb-go/chat/protos"
	"chat/settings"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
	"time"
)

var (
	address  = fmt.Sprintf("%s:%s", settings.ChatHost, settings.ChatPort)
	command  string
	login    string
	password string
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к %s", address)
	}
	defer conn.Close()

	fmt.Println("Приложение чата по gRPC")
	client := pb.NewChatClient(conn)

	//auth(client)
	login = "rizhiy"
	fmt.Printf("Добро пожаловать, %s! Вот список доступных команд:\n", login)
	fmt.Println("/open USERNAME - открыть диалог с пользователем")
	fmt.Println("/send USERNAME MESSAGE - отправить сообщение пользователю")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nВведите команду: ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			return
		}

		command = strings.TrimSuffix(command, "\n")
		commandLine := strings.Fields(command)
		action := commandLine[0]
		interlocutor := commandLine[1]

		if action == "/send" {
			message := strings.Join(commandLine[2:], " ")
			currTime := time.Now().In(time.Now().Location()).Unix()
			send(client, login, interlocutor, message, currTime)
		} else if action == "/open" {
			get(client, login, interlocutor)
		} else {
			fmt.Println("Не удалось распознать команду")
		}
	}
}

// auth позволяет провести аутентификацию пользователя
func auth(c pb.ChatClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	fmt.Println("Необходимо пройти аутентификацию")

	for {
		fmt.Print("Введите логин: ")
		_, err := fmt.Scanln(&login)
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			return
		}

		fmt.Print("Введите пароль: ")
		_, err = fmt.Scanln(&password)
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			return
		}

		isAuth, err := c.AuthUser(ctx, &pb.ChatAuth{Login: login, Password: password})
		if err != nil {
			fmt.Println("Ошибка при вызове AuthUser:", err)
			return
		}

		if !isAuth.Errors {
			fmt.Printf("Аутентификация пройдена\n\n")
			break
		} else {
			fmt.Printf("Аутентификация не пройдена\n\n")
			continue
		}
	}
}

// get позволяет получить сообщения из чата с пользователем
func get(c pb.ChatClient, user, interlocutor string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	messages, err := c.UpdateChat(ctx, &pb.ChatUpdate{From: user, To: interlocutor})
	if err != nil {
		log.Fatalf("Ошибка при вызове UpdateChat: %v", err)
	}

	for _, msg := range messages.Message {
		fmt.Printf("От %s в %d: %s\n", msg.From, msg.Timestamp, msg.Message)
	}
}

// send позволяет отправить данные пользователю
func send(c pb.ChatClient, user, interlocutor, message string, timestamp int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	sending, err := c.SaveMessage(ctx, &pb.ChatMessage{Timestamp: timestamp, From: user, To: interlocutor, Message: message})
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	if sending.Errors {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}
}
