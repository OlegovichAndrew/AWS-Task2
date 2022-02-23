package main

import (
	"flag"
	"log"
	"net"

	"aws-server/proto"
	"aws-server/transport"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "The server address")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterAWSServiceServer(s, &transport.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
