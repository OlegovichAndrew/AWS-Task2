package transport

import (
	"aws-server/proto"
	"context"
	"log"
)

type Server struct {
	proto.UnimplementedAWSServiceServer
}

func (s *Server) StringSend(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	log.Printf("Received: %v", request.GetTestReq())
	return &proto.Response{TestResp: "Request received:" + request.GetTestReq()}, nil
}
