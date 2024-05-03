package main

import (
	"chat/database"
	"chat/settings"

	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "chat/gen/pb-go/chat/protos"
)

type chatServer struct {
	pb.UnimplementedChatServer
}

// AuthUser это функция проверки наличия пользователя в базе данных
func (s *chatServer) AuthUser(ctx context.Context, auth *pb.ChatAuth) (*pb.ChatResponse, error) {
	flag, err := database.AuthUser(auth.Login, auth.Password)

	if err != nil {
		return &pb.ChatResponse{
			Errors:    true,
			ErrorCode: fmt.Sprintf("Ошибка на стороне базы данных при создании пользователя %s\n", auth.Login),
		}, nil
	}

	if !flag {
		return &pb.ChatResponse{
			Errors:    true,
			ErrorCode: "Неверный логин или пароль",
		}, nil
	}

	return &pb.ChatResponse{
		Errors:    false,
		ErrorCode: "",
	}, nil
}

// CreateUser это функция создания пользователя в базе данных
func (s *chatServer) CreateUser(ctx context.Context, auth *pb.ChatAuth) (*pb.ChatResponse, error) {
	// Здесь можно реализовать логику регистрации
	err := database.CreateUser(auth.Login, auth.Password)
	if err != nil {
		return &pb.ChatResponse{
			Errors:    true,
			ErrorCode: fmt.Sprintf("Не удалось создать пользователя %s\n", auth.Login),
		}, nil
	}

	return &pb.ChatResponse{
		Errors:    false,
		ErrorCode: "",
	}, nil
}

// SaveMessage это функция сохранения сообщения в базе данных
func (s *chatServer) SaveMessage(ctx context.Context, message *pb.ChatMessage) (*pb.ChatResponse, error) {
	err := database.CreateMessage(message.Timestamp, message.From, message.To, message.Message)
	if err != nil {
		return &pb.ChatResponse{
			Errors:    true,
			ErrorCode: fmt.Sprintf("%d Новое сообщение от %s для %s: %s", message.Timestamp, message.From, message.To, message.Message),
		}, nil
	}

	return &pb.ChatResponse{
		Errors:    false,
		ErrorCode: "",
	}, nil
}

// UpdateChat это функция обновления сообщений в чате
func (s *chatServer) UpdateChat(ctx context.Context, update *pb.ChatUpdate) (*pb.ChatMessages, error) {
	messages, err := database.GetMessages(update.From, update.To)
	if err != nil {
		return nil, err
	}

	return &pb.ChatMessages{
		Message: messages,
	}, nil
}

func serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.ChatHost, settings.ChatPort))
	if err != nil {
		log.Fatalf("Ошибка инициализации сокета: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &chatServer{})
	reflection.Register(s)

	log.Printf(fmt.Sprintf("Старт gRPC сервера %s:%s", settings.ChatHost, settings.ChatPort))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func main() {
	database.Migration()
	serve()
}
