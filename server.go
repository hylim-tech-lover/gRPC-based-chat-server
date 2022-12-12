package main

import (
	"log"
	"net"
	"os"

	chat_server "grpcChatServer/chat_server"

	"google.golang.org/grpc"
)

func main() {
	// assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		// Set default port (8080)
		Port = "8080"
	}
	// init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen at %v :: %v. Server will be terminated", Port, err)
	}
	log.Printf("Listening at %v", Port)

	// gRPC server instance
	grpcServer := grpc.NewServer()

	// Register ChatService
	cs := chat_server.ChatServer{}

	chat_server.RegisterServicesServer(grpcServer, &cs)

	// gRPC listen and serve
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC with error %v.Server will be terminated", err)
	}

}
