package main

import (
	"testing"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
)

type FakeServer struct{}

func (f *FakeServer) callServiceBusPut(in *pb.Request) ([]byte, error) {
	return []byte{1, 2, 2}, nil
}

func TestGetReleaseTagMessage(t *testing.T) {
	f := &FakeServer{}
	_, err := f.callServiceBusPut(&pb.Request{})
	if err != nil {
		t.Fatalf("Expected err to be nil but it was %s", err)
	}
}
