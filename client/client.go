package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "", "The server name use to verify the hostname returned by TLS handshake")
	contractorID       = flag.Int64("contractorid", 72494, "Contractor Id for the PUT call")
	orderNumber        = flag.Int64("ordernumber", 600016555, "OrderNumber for the PUT call")
	imageType          = flag.Int("imagetype", 1, "Imagetype for the PUT call")
	fileName           = flag.String("filename", "../testdata/e3e0f976-79a5-4059-ac23-d44386a6d4da.png", "Filename for the PUT call")
	imageWidth         = flag.Int("imagewidth", 100, "Imagewidth for the PUT call")
	imageHeight        = flag.Int("imageheight", 100, "Imageheight for the PUT call")
	releaseDate        = flag.String("releasedate", "2015-08-06", "Releasedate for the PUT call")
	deptCode           = flag.String("deptcode", "01", "Department code for the PUT call")
)

type input struct {
	contracttorid int64
	ordernumber   int64
	imagetype     int32
	filename      string
	imagewidth    int32
	imageheight   int32
	releasedate   string
	deptcode      string
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		var sn string
		if *serverHostOverride != "" {
			sn = *serverHostOverride
		}
		var creds credentials.TransportCredentials
		if *caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewContentServiceClient(conn)

	// Create PutRequest
	in := &input{}
	in.contracttorid = *contractorID
	in.deptcode = *deptCode
	in.filename = *fileName
	in.imageheight = int32(*imageHeight)
	in.imagetype = int32(*imageType)
	in.imagewidth = int32(*imageWidth)
	in.ordernumber = *orderNumber
	in.releasedate = *releaseDate
	putRequest, err := createPutRequest(in)
	if err != nil {
		log.Fatalf("Error creating put request: %v", err)
	}

	// Contact the server and print out its response.
	response, err := client.Put(context.Background(), putRequest)
	if err != nil {
		log.Fatalf("Error making put call: %v", err)
	}
	if response.GetResult() != nil {
		fmt.Println("put id: ", response.GetResult().GetId())
	} else {
		fmt.Println("put error: ", response.GetError().GetMessage())
	}
}

func createPutRequest(in *input) (*pb.PutRequest, error) {
	putRequest := &pb.PutRequest{}
	putRequest.Contractorid = in.contracttorid
	putRequest.Deptcode = in.deptcode
	putRequest.Filename = in.filename
	putRequest.Imageheight = in.imageheight
	putRequest.Imagetype = in.imagetype
	putRequest.Imagewidth = in.imagewidth
	putRequest.Ordernumber = in.ordernumber
	putRequest.Releasedate = in.releasedate
	fileContents, err := getFileContents(in.filename)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving the file contents: %s", err)
	}
	putRequest.Filecontents = fileContents

	return putRequest, nil
}

func getFileContents(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := stats.Size()
	fileBytes := make([]byte, size)
	reader := bufio.NewReader(file)
	_, err = reader.Read(fileBytes)
	if err != nil {
		return "", err
	}
	encodedFileBytes := base64.StdEncoding.EncodeToString(fileBytes)

	return encodedFileBytes, nil
}
