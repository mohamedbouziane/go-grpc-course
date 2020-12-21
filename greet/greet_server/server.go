package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Printf("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error : %v", err)
	}

}
