package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

type FakeServer struct {
	Response *pb.JSONRPCResponse
	Err      error
}

func (f *FakeServer) callServiceBus(request *pb.JSONRPCRequest) (*pb.JSONRPCResponse, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Response, nil
}

var cases = []struct {
	server           *Server
	fakeServer       *FakeServer
	request          *pb.PutRequest
	expectedResponse *pb.PutResponse
	expectedErr      error
}{
	{
		server: NewServer("http://servicebus.qa01.local/Execute.svc/Execute"),
		fakeServer: &FakeServer{
			Response: &pb.JSONRPCResponse{Jsonrpc: "2.0",
				Result: &pb.JSONRPCResult{Contractorid: 72494,
					Releasedate:   "2015-08-06T15:09:30",
					Scandate:      "2017-03-09T10:33:09",
					Imagetype:     1,
					Imagewidth:    100,
					Imageheight:   100,
					Deptcode:      "01",
					Descprefix:    "test",
					Desctext:      "test",
					Category:      "test",
					Ordernumber:   600016555,
					Archived:      "N",
					Datecreated:   "2017-03-09T10:33:09",
					Datemodefied:  "2017-03-09T10:33:09",
					Filesize:      180,
					Id:            1810448062,
					Imagefilename: "\\\\filer\\QA01\\ImageStore\\ServiceBus\\600\\016\\555\\da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
					Thumbnailsize: 0,
					Webfilename:   "QA01/ImageStore/ServiceBus/600/016/555/da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
				}},
			Err: nil,
		},
		request: &pb.PutRequest{Contractorid: 72494,
			Ordernumber: 600016555,
			Imagetype:   1,
			Filename:    "test.png",
			Imagewidth:  100,
			Imageheight: 100,
			Releasedate: "2015-08-06",
			Deptcode:    "01",
		},
		expectedResponse: &pb.PutResponse{Result: &pb.JSONRPCResult{Contractorid: 72494,
			Releasedate:   "2015-08-06T15:09:30",
			Scandate:      "2017-03-09T10:33:09",
			Imagetype:     1,
			Imagewidth:    100,
			Imageheight:   100,
			Deptcode:      "01",
			Descprefix:    "test",
			Desctext:      "test",
			Category:      "test",
			Ordernumber:   600016555,
			Archived:      "N",
			Datecreated:   "2017-03-09T10:33:09",
			Datemodefied:  "2017-03-09T10:33:09",
			Filesize:      180,
			Id:            1810448062,
			Imagefilename: "\\\\filer\\QA01\\ImageStore\\ServiceBus\\600\\016\\555\\da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
			Thumbnailsize: 0,
			Webfilename:   "QA01/ImageStore/ServiceBus/600/016/555/da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
		}},
		expectedErr: nil,
	},
	{
		server: NewServer("http://servicebus.qa01.local/Execute.svc/Execute"),
		fakeServer: &FakeServer{
			Response: &pb.JSONRPCResponse{},
			Err:      errors.New("Fake Error"),
		},
		request:          &pb.PutRequest{},
		expectedResponse: nil,
		expectedErr:      errors.New("Fake Error"),
	},
}

var jsonRPCRequestCases = []struct {
	putRequest          *pb.PutRequest
	expectedJSONRequest *pb.JSONRPCRequest
}{
	{
		putRequest: &pb.PutRequest{Contractorid: 72494,
			Ordernumber: 600016555,
			Imagetype:   1,
			Filename:    "test.png",
			Imagewidth:  100,
			Imageheight: 100,
			Releasedate: "2015-08-06",
			Deptcode:    "01",
		},
		expectedJSONRequest: &pb.JSONRPCRequest{Jsonrpc: "2.0",
			Method: "CONTENTSERVICE.PUT",
			Params: &pb.PutRequest{Contractorid: 72494,
				Ordernumber: 600016555,
				Imagetype:   1,
				Filename:    "test.png",
				Imagewidth:  100,
				Imageheight: 100,
				Releasedate: "2015-08-06",
				Deptcode:    "01",
			}},
	},
}

var putResponseCases = []struct {
	jsonResponse        *pb.JSONRPCResponse
	expectedPutResponse *pb.PutResponse
}{
	{
		jsonResponse: &pb.JSONRPCResponse{Jsonrpc: "2.0",
			Result: &pb.JSONRPCResult{Contractorid: 72494,
				Releasedate:   "2015-08-06T15:09:30",
				Scandate:      "2017-03-09T10:33:09",
				Imagetype:     1,
				Imagewidth:    100,
				Imageheight:   100,
				Deptcode:      "01",
				Descprefix:    "test",
				Desctext:      "test",
				Category:      "test",
				Ordernumber:   600016555,
				Archived:      "N",
				Datecreated:   "2017-03-09T10:33:09",
				Datemodefied:  "2017-03-09T10:33:09",
				Filesize:      180,
				Id:            1810448062,
				Imagefilename: "\\\\filer\\QA01\\ImageStore\\ServiceBus\\600\\016\\555\\da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
				Thumbnailsize: 0,
				Webfilename:   "QA01/ImageStore/ServiceBus/600/016/555/da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
			}},
		expectedPutResponse: &pb.PutResponse{Result: &pb.JSONRPCResult{Contractorid: 72494,
			Releasedate:   "2015-08-06T15:09:30",
			Scandate:      "2017-03-09T10:33:09",
			Imagetype:     1,
			Imagewidth:    100,
			Imageheight:   100,
			Deptcode:      "01",
			Descprefix:    "test",
			Desctext:      "test",
			Category:      "test",
			Ordernumber:   600016555,
			Archived:      "N",
			Datecreated:   "2017-03-09T10:33:09",
			Datemodefied:  "2017-03-09T10:33:09",
			Filesize:      180,
			Id:            1810448062,
			Imagefilename: "\\\\filer\\QA01\\ImageStore\\ServiceBus\\600\\016\\555\\da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
			Thumbnailsize: 0,
			Webfilename:   "QA01/ImageStore/ServiceBus/600/016/555/da00563b-bb38-49b1-b3ef-29dbce63fbed.png",
		}},
	},
}

func TestCallServiceBus(t *testing.T) {
	for _, c := range cases {
		c.server.ServiceBusCaller = c.fakeServer
		response, err := c.server.Put(context.Background(), c.request)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

		if !reflect.DeepEqual(c.expectedResponse, response) {
			t.Errorf("Expected %q but got %q", c.expectedResponse, response)
		}
	}
}

func TestCreateJSONRPCRequest(t *testing.T) {
	for _, c := range jsonRPCRequestCases {
		jsonRequest := createJSONRPCRequest(c.putRequest)
		if !reflect.DeepEqual(jsonRequest, c.expectedJSONRequest) {
			t.Errorf("Expected %q but got %q", c.expectedJSONRequest, jsonRequest)
		}
	}
}

func TestCreatePutResponse(t *testing.T) {
	for _, c := range putResponseCases {
		putResponse := createPutResponse(c.jsonResponse)
		if !reflect.DeepEqual(c.expectedPutResponse, putResponse) {
			t.Errorf("Expected %q but got %q", c.expectedPutResponse, putResponse)
		}
	}
}
