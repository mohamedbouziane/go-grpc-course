package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/go-grpc-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'am a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)
	doServerStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do a Unary RPC...")
	req := &greetpb.GreetingRequest{

		Greeting: &greetpb.Greeting{
			FirstName: "Mohamed",
			LastName:  "Bouziane",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC : %v", err)
	}
	log.Panicf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{

		Greeting: &greetpb.Greeting{
			FirstName: "Mohamed",
			LastName:  "Bouziane",
		},
	}
   
	reaStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while caling GreetManyTimes : %v", err)
	}

	for {

		msg, err := reaStream.Recv()
		if err == io.EOF{
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading the stream: %v", err)
		}
		log.Printf("Response fron the GreetManyTimes; %v", msg.GetResult())
	}


}