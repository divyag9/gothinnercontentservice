package main

import (
	"errors"
	"reflect"
	"testing"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

type FakeServer struct {
	Response *pb.JSONRPCResponse
	Err      error
}

func (f *FakeServer) CallServiceBusPut(request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

var cases = []struct {
	f                FakeServer
	request          *pb.JSONRPCRequest
	expectedResponse *pb.JSONRPCResponse
	expectedErr      error
}{
	{
		f: FakeServer{
			Response: &pb.JSONRPCResponse{},
			Err:      nil,
		},
		request: &pb.JSONRPCRequest{Jsonrpc: "2.0",
			Method: "CONTENTSERVICE.PUT",
			Params: &pb.PutRequest{Contractorid: 72494,
				Ordernumber: 600016555,
				Imagetype:   1,
				Filename:    "test.png",
				Imagewidth:  100,
				Imageheight: 100,
				Releasedate: "2015-08-06",
				Deptcode:    "01",
			},
		},
		expectedResponse: &pb.JSONRPCResponse{},
		expectedErr:      nil,
	},
	{
		f: FakeServer{
			Response: &pb.JSONRPCResponse{},
			Err:      errors.New("Fake Error"),
		},
		request:          &pb.JSONRPCRequest{},
		expectedResponse: nil,
		expectedErr:      errors.New("Fake Error"),
	},
}

func TestCallServiceBusPut(t *testing.T) {
	for _, c := range cases {
		response, err := c.f.CallServiceBusPut(c.request)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

		if !reflect.DeepEqual(c.expectedResponse, response) {
			t.Errorf("Expected %q but got %q", c.expectedResponse, response)
		}
	}
}
