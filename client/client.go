package main

import (
	"context"
	"log"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewContentServiceClient(conn)

	// Contact the server and print out its response.
	_, err = c.Put(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("error making put call: %v", err)
	}
}
