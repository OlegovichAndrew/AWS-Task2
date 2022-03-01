package main

import (
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
	addr            = flag.String("addr", "localhost:4444", "the address to connect to")
)

func main() {
	flag.StringVar(&bucket, "b", "upload.practice", "The bucket to download/upload the file from/to")
	flag.StringVar(&fileKey, "f", "The_first_upload/number.txt", "The file to download/upload")
	flag.Parse()

	if bucket == "" || fileKey == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
		return
	}

	lis, err := net.Listen("tcp", *addr)
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
