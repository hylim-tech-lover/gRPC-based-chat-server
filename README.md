# Bidirectional Chat Server using gRPC in Go

### Tech stack:

- gRPC server (`protobuf` to `go` file)
- Go

### General workflow:

> Following build steps are just for initial setup and not needed to execute as the steps have been executed and application can be executed directly using Go.

1. Generate `go.mod` via `go mod init grpcChatServer`
2. Create `chat.proto` file and `chat_server` folder
3. Install following dependencies based on OS:
   - protobuf
   - protoc-gen-go
   - protoc-gen-go-grpc
4. Generate following in `chat_server` folder from `chat.proto` file via `protoc --go-grpc_out=require_unimplemented_servers=false:./chat_server/ --go_out=./chat_server/ chat.proto`:
   - chat_grpc.pb.go
   - chat.pb.go
5. Install `google.golang.org/grpc` using `go get` command.
6. Create `chat_server.go` as main server function to forward and route client message.
7. Create `server.go` as gRPC server.
8. Create `client.go` as client that connected to gRPC server and send and receive message.
9. Application at client and server side will terminate once connection is disconnected.

---

## Installation

1. Before start, hosting machine should be installed with [Go](https://go.dev/doc/install)

2. Data Structure - [chat.proto](./chat.proto)

---

## Execution

1. Git clone this repository to local and open with your preferred IDE.
2. Open a terminal and execute `go run server.go`. Alternatively, execute `env PORT=9000 go run server.go` to change port number to `9000`. Default port is `8080`.
3. Open another terminal and execute `go run client.go`
4. Key in host address [`localhost`] with corresponding port stated in Step 2.
5. Key in username and start sending message.
6. Repeat step 3 and 5 for another clients or more.
7. All clients should received the message sent by other clients which simulated the chat server.
