package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

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

func (*server) 	GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error{
	fmt.Printf("GreetManyTimes function was invoked with %v", req)
	
	firstName := req.GetGreeting().GetFirstName()
	for i:=0; i<=10; i++{
		result := "Hello " + firstName +" number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Microsecond)	

	}

	return nil

}


func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error{
	fmt.Printf("LongGreet function was invoked with a streaming request...\n")
	result :=""

	for{
		req, err := stream.Recv()
		if err== io.EOF{
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil{
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result+= "Hello " + firstName + "! "

	}

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
