package main

import (
	"aws-dl-s3/config"
	"aws-dl-s3/transport"
	"flag"
	"fmt"
	"log"
	"net"

	"aws-dl-s3/proto"
	"google.golang.org/grpc"
)

var (
	bucket, fileKey string
)

func main() {
	flag.StringVar(&bucket, "b", "ul.practice", "The bucket to download/upload the file from/to")
	flag.StringVar(&fileKey, "f", "number.txt", "The file to download/upload")
	flag.Parse()

	if bucket == "" || fileKey == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
		return
	}

	lis, err := net.Listen("tcp", config.GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := transport.NewServer()
	s := grpc.NewServer()
	proto.RegisterAWSServiceServer(s, server)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
