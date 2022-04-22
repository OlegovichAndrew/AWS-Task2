package main

import (
	"aws-dl-s3/config"
	"aws-dl-s3/transport"
	"log"
	"net"

	"aws-dl-s3/proto"
	"google.golang.org/grpc"
)

func main() {

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
