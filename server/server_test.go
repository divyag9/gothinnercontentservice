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
		request:          &pb.JSONRPCRequest{},
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

var encodeJSONRPCRequestCases = []struct {
	putRequest          *pb.PutRequest
	expectedJSONRequest *pb.JSONRPCRequest
	expectedErr         error
}{
	{
		putRequest:          &pb.PutRequest{},
		expectedJSONRequest: &pb.JSONRPCRequest{},
		expectedErr:         nil,
	},
	{
		putRequest:          &pb.PutRequest{},
		expectedJSONRequest: nil,
		expectedErr:         errors.New("Fake Error"),
	},
}

var encodePutResponseCases = []struct {
	jsonResponse        *pb.JSONRPCResponse
	expectedPutResponse *pb.PutResponse
	expectedErr         error
}{
	{
		jsonResponse:        &pb.JSONRPCResponse{},
		expectedPutResponse: &pb.PutResponse{},
		expectedErr:         nil,
	},
	{
		jsonResponse:        &pb.JSONRPCResponse{},
		expectedPutResponse: nil,
		expectedErr:         errors.New("Fake Error"),
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

func TestEncodeJSONRPCRequest(t *testing.T) {
	s := &server{}
	for _, c := range encodeJSONRPCRequestCases {
		jsonRequest, err := s.EncodeJSONRPCRequest(c.putRequest)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

		if !reflect.DeepEqual(c.expectedJSONRequest, jsonRequest) {
			t.Errorf("Expected %q but got %q", c.expectedJSONRequest, jsonRequest)
		}
	}
}

func TestEncodePutResponse(t *testing.T) {
	s := &server{}
	for _, c := range encodePutResponseCases {
		putResponse, err := s.EncodePutResponse(c.jsonResponse)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

		if !reflect.DeepEqual(c.expectedPutResponse, putResponse) {
			t.Errorf("Expected %q but got %q", c.expectedPutResponse, putResponse)
		}
	}
}
