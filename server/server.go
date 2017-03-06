package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

type server struct{}

func (s *server) Put(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	//encodeJSON
	//callServiceBusPut
	//decodeJSON
	return &pb.Response{}, nil
}

func (s *server) callServiceBusPut(in *pb.Request) ([]byte, error) {
	return []byte{1, 2, 2}, nil
}

func encodeJSON(in *pb.Request) ([]byte, error) {
	return []byte{1, 2, 2}, nil
}

func decodeJSON(result []byte) (*pb.Response, error) {
	return &pb.Response{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterContentServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
