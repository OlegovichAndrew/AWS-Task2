package main

import (
	"log"
	"net"

	"aws-grpc-server/config"
	"aws-grpc-server/proto"
	"aws-grpc-server/transport"

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
