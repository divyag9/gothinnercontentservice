package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile           = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile            = flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port               = flag.Int("port", 10000, "The server port")
	serviceBusEndPoint = flag.String("servicebus_endpoint", "http://servicebus.qa01.local/Execute.svc/Execute", "The servicebus execute endpoint")
)

type server struct {
	serviceBusEndPoint string
}

func (s *server) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	jsonRPCRequest := createJSONRPCRequest(request)
	jsonRPCResponse, err := getServiceBusResponse(s, jsonRPCRequest)
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

func getServiceBusResponse(s *server, request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error) {
	requestBytes, err := marshalJSONRPCRequest(request)
	if err != nil {
		return nil, err
	}
	responseBody, err := s.CallServiceBusPut(requestBytes)
	if err != nil {
		return nil, err
	}
	jsonRPCResponse := &pb.JSONRPCResponse{}
	body, _ := ioutil.ReadAll(responseBody)
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	jsonRPCResponse, err = unmarshalJSONRPCResponse(body)
	fmt.Println("json response: ", jsonRPCResponse.Result)
	if err != nil {
		return nil, err
	}

	return jsonRPCResponse, nil
}

func marshalJSONRPCRequest(request *pb.JSONRPCRequest) ([]byte, error) {
	jsonRPCRequestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return jsonRPCRequestBytes, nil
}

func (s *server) CallServiceBusPut(requestBytes []byte) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", s.serviceBusEndPoint, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func unmarshalJSONRPCResponse(result []byte) (*pb.JSONRPCResponse, error) {
	jsonRPCResponse := &pb.JSONRPCResponse{}
	err := json.Unmarshal(result, jsonRPCResponse)
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

func main() {
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
	contentServiceServer := &server{serviceBusEndPoint: *serviceBusEndPoint}
	pb.RegisterContentServiceServer(grpcServer, contentServiceServer)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}
