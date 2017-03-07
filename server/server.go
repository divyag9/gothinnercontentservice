package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type server struct{}

func (s *server) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	jsonRPCRequest, err := s.EncodeJSONRPCRequest(request)
	if err != nil {
		return nil, err
	}
	jsonRPCResponse, err := s.CallServiceBusPut(jsonRPCRequest)
	if err != nil {
		return nil, err
	}
	putResponse, err := s.EncodePutResponse(jsonRPCResponse)
	if err != nil {
		return nil, err
	}
	return putResponse, nil
}

func (s *server) EncodeJSONRPCRequest(request *pb.PutRequest) (*pb.JSONRPCRequest, error) {
	return &pb.JSONRPCRequest{}, nil
}

func (s *server) CallServiceBusPut(request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error) {
	return &pb.JSONRPCResponse{}, nil
}

func marshalJSONRPCRequest(request *pb.JSONRPCRequest) ([]byte, error) {
	return []byte{1, 2, 2}, nil
}

func unmarshalJSONRPCResponse(result []byte) (*pb.JSONRPCResponse, error) {
	return &pb.JSONRPCResponse{}, nil
}

func (s *server) EncodePutResponse(response *pb.JSONRPCResponse) (*pb.PutResponse, error) {
	return &pb.PutResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterContentServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err) // We want to continue serving and not die?
	}
}
