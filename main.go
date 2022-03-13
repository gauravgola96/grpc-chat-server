package main

import (
	"google.golang.org/grpc"
	"grpc-chat-server/proto"
	"log"
	"net"
	"os"
)

func main() {

	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5005" //default Port set to 5000 if PORT is not set in env
	}

	//init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := ChatServer{}
	proto.RegisterServicesServer(grpcserver, &cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}

}
