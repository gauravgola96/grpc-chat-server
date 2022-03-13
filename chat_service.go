package main

import (
	"bufio"
	"grpc-chat-server/proto"
	"log"
	"os"
	"strings"
)

type ChatServer struct {
	proto.UnimplementedServicesServer
}

func (is *ChatServer) ChatService(srv proto.Services_ChatServiceServer) error {

	errch := make(chan error)

	// receive messages - init a go routine
	go receiveFromStream(srv, errch)

	// send messages - init a go routine
	go sendToStream(srv, errch)

	return <-errch

}

func receiveFromStream(srv proto.Services_ChatServiceServer, errch_ chan error) {

	for {
		msg, err := srv.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch_ <- err
			return
		}
		log.Printf("received  %s", msg.Body)
	}
}

func sendToStream(srv proto.Services_ChatServiceServer, errch_ chan error) {

	for {
		//fmt.Printf("Type your thing : ")
		reader := bufio.NewReader(os.Stdin)
		serverMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		serverMessage = strings.Trim(serverMessage, "\r\n")

		clientMessageBox := &proto.Response{
			Name: "bot",
			Body: serverMessage,
		}

		err = srv.Send(clientMessageBox)
		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}
