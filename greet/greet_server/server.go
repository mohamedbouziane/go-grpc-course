package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-grpc-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	fristName := req.GetGreeting().GetFirstName()
	result := "Hello " + fristName
	res := &greetpb.GreetingResponse{
		Result: result,
	}
	return res, nil
}

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
