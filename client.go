package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	chat_server "grpcChatServer/chat_server"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Enter Server IP:Port :: [Eg: 192.168.0.1:3000]")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console : %v ", err)
	}
	serverID = strings.Trim(serverID, "\r\n")
	log.Printf("Connecting to %v ", serverID)

	// Connect to grpc server - not using SSL/TLS encryption
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server : %v", err)
	}

	defer conn.Close()

	// Call ChatService to create a stream
	client := chat_server.NewServicesClient(conn)

	// Passed in a context with no deadline
	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService : %v", err)
	}

	log.Printf("%s Connected!", serverID)

	// Implement communication with gRPC server
	ch := clientHandle{stream: stream}
	ch.clientConfig()

	// init go routine to asynchronously run
	go ch.sendMessage()
	go ch.receiveMessage()

	// Blocker channel
	bl := make(chan bool)
	<-bl
}

// Client handle struct
type clientHandle struct {
	stream     chat_server.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	// Ask username
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please insert your display name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read from console : %v ", err)
	}

	ch.clientName = strings.Trim(name, "\r\n")
}

// Send message
func (ch *clientHandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("[User: %s ] Insert message: ", ch.clientName)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read from console : %v ", err)
		}

		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chat_server.ClientMsg{Name: ch.clientName, Body: clientMessage}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server : %v", err)
		}
	}
}

// Receive message
func (ch *clientHandle) receiveMessage() {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server : %v", err)
		}

		// Print message to console
		fmt.Printf("\n %s : %s \n", msg.Name, msg.Body)
	}
}
