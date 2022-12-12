package __

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct{}

// Define ChatService
func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// Receive message - init a go routine to asynchronously run
	go receiveFromStream(csi, clientUniqueCode, errch)

	// Sending message - init a go routine to asynchronously run
	go sendToStream(csi, clientUniqueCode, errch)

	// Terminate connection when receiving error through channel errch
	return <-errch
}

// Receive message function
func receiveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	// Implement a loop that continously receive message from client
	for {
		msg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client: %v", err)
			errch_ <- err
		} else {
			// Add message to queue
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        msg.Name,
				MessageBody:       msg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode_,
			})

			messageHandleObject.mu.Unlock()

			// Show latest message of the queue
			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}

	}
}

// Send message
func sendToStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	// Implement a loop
	for {
		// Loop through message in queue
		for {
			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName := messageHandleObject.MQue[0].ClientName
			senderMessage := messageHandleObject.MQue[0].MessageBody

			// Send message to designated client (do not send to the same client)
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&ServerMsg{Name: senderName, Body: senderMessage})

				if err != nil {
					errch_ <- err
				}

				if len(messageHandleObject.MQue) > 1 {
					// Delete first message once sent to the designated client
					messageHandleObject.MQue = messageHandleObject.MQue[1:]
				} else {
					// Clear whole queue if it is the last message in the queue
					messageHandleObject.MQue = []messageUnit{}
				}
			}
			messageHandleObject.mu.Unlock()
		}
	}

}
