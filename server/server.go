package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

// ServiceBusCaller interface for calling ServiceBus
type ServiceBusCaller interface {
	callServiceBus(request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error)
}

// Server servicebus
type Server struct {
	ServiceBusCaller
}

// Caller interface for servicebus
type Caller struct {
	serviceBusEndPoint string
}

// Put performs the servicebus put
func (s *Server) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	jsonRPCRequest := createJSONRPCRequest(request)
	jsonRPCResponse, err := s.ServiceBusCaller.callServiceBus(jsonRPCRequest)
	if err != nil {
		return nil, err
	}
	putResponse := createPutResponse(jsonRPCResponse)

	return putResponse, nil
}

func createJSONRPCRequest(request *pb.PutRequest) *pb.JSONRPCRequest {
	jsonRPCRequest := &pb.JSONRPCRequest{}
	jsonRPCRequest.Jsonrpc = "2.0"
	jsonRPCRequest.Method = "CONTENTSERVICE.PUT"
	jsonRPCRequest.Params = request

	return jsonRPCRequest
}

func (c *Caller) callServiceBus(request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error) {
	start := time.Now()
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed Marshal: ", elapsed)

	req, err := http.NewRequest("POST", c.serviceBusEndPoint, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	jsonRPCResponse := &pb.JSONRPCResponse{}
	err = json.Unmarshal(body, jsonRPCResponse)
	if err != nil {
		return nil, err
	}

	return jsonRPCResponse, nil
}

func createPutResponse(response *pb.JSONRPCResponse) *pb.PutResponse {
	putResponse := &pb.PutResponse{}
	putResponse.Result = response.GetResult()
	putResponse.Error = response.GetError()

	return putResponse
}

// NewServer creates new servicebus  server
func NewServer(serviceBusEndPoint string) *Server {
	return &Server{ServiceBusCaller: &Caller{serviceBusEndPoint: serviceBusEndPoint}}
}

func main() {
	tls := flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile := flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile := flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port := flag.Int("port", 10000, "The server port")
	serviceBusEndPoint := flag.String("servicebus_endpoint", "http://servicebus.qa01.local/Execute.svc/Execute", "The servicebus execute endpoint")

	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
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
	server := NewServer(*serviceBusEndPoint)
	pb.RegisterContentServiceServer(grpcServer, server)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}
