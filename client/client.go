package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/divyag9/gothinnercontentservice/contentservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "test.com", "The server name use to verify the hostname returned by TLS handshake")
	contractorID       = flag.Int64("contractor_id", 0, "Contractor Id for the PUT call")
	orderNumber        = flag.Int64("order_number", 0, "OrderNumber for the PUT call")
	imageType          = flag.Int("image_type", 0, "Imagetype for the PUT call")
	fileName           = flag.String("file_name", "", "Filename for the PUT call")
	imageWidth         = flag.Int("image_width", 0, "Imagewidth for the PUT call")
	imageHeight        = flag.Int("image_height", 0, "Imageheight for the PUT call")
	releaseDate        = flag.String("release_date", "", "Releasedate for the PUT call")
	deptCode           = flag.String("dept_code", "", "Department code for the PUT call")
)

type input struct {
	contracttorid int64
	ordernumber   int64
	imagetype     int
	filename      string
	imagewidth    int
	imageheight   int
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
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewContentServiceClient(conn)

	in := &input{}
	in.contracttorid = *contractorID
	in.deptcode = *deptCode
	in.filename = *fileName
	in.imageheight = *imageHeight
	in.imagetype = *imageType
	in.imagewidth = *imageWidth
	in.ordernumber = *orderNumber
	in.releasedate = *releaseDate

	putRequest, err := createPutRequest(in)
	if err != nil {
		log.Fatalf("error creating put request: %v", err)
	}
	// Contact the server and print out its response.
	_, err = client.Put(context.Background(), putRequest)
	if err != nil {
		log.Fatalf("error making put call: %v", err)
	}

}

func createPutRequest(in *input) (*pb.PutRequest, error) {
	//get the filecontents base64 encode and set to putrequest filecontents
	return &pb.PutRequest{}, nil
}
