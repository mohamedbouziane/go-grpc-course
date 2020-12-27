package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/go-grpc-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)

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
	fmt.Println("Starting to do a Server Streaming RPC...")
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
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading the stream: %v", err)
		}
		log.Printf("Response fron the GreetManyTimes; %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{

		&greetpb.LongGreetRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Mohamed",
			},
		},

		&greetpb.LongGreetRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Ismail",
			},
		},

		&greetpb.LongGreetRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Anes",
			},
		},

		&greetpb.LongGreetRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Yahia",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet; %v", err)
	}

	// we iterate over our slice and send each message individualy

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Microsecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream bu invoking the client

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{

		&greetpb.GreetEveryoneRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Mohamed",
			},
		},

		&greetpb.GreetEveryoneRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Ismail",
			},
		},

		&greetpb.GreetEveryoneRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Anes",
			},
		},

		&greetpb.GreetEveryoneRequest{

			Greeting: &greetpb.Greeting{
				FirstName: "Yahia",
			},
		},
	}
	waitc := make(chan struct{})
	// we send a bunch of message to the client (go routine)
	go func() {
		// function to send a messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Microsecond)
		}
		stream.CloseSend()

	}()

	// we receive a bunch of messages from the client (go routine)

	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}

			fmt.Printf("Receiving: %v \n", res.GetResult())
		}
		close(waitc)

	}()
	// block untile everything is done
	<-waitc

}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohamed",
			LastName:  "Bouziane",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
