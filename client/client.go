package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-chat-server/proto"
	"log"
	"os"
	"strings"
)

func main() {
	serverID := "localhost:5005"

	log.Println("Connecting : " + serverID)

	//connect to grpc server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	//call ChatService to create a stream
	client := proto.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clienthandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl

}

//clienthandle
type clienthandle struct {
	stream     proto.Services_ChatServiceClient
	clientName string
}

func (ch *clienthandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")

}

//send message
func (ch *clienthandle) sendMessage() {

	// create a loop
	fmt.Printf("Type your thing : ")
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &proto.Request{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

//receive message
func (ch *clienthandle) receiveMessage() {

	//create a loop
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		fmt.Printf("received : %s \n", msg.Body)

	}
}
