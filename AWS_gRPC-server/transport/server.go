package transport

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"aws-grpc-server/proto"
	"aws-grpc-server/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc/peer"
)

var client *s3.Client

type Server struct {
	proto.UnimplementedAWSServiceServer
}

func NewServer() *Server {
	return &Server{}
}

//Download function downloads the file from the bucket came with request. Then send it back by grpc steram.
func (s *Server) Download(req *proto.Request, responseStream proto.AWSService_DownloadServer) error {
	// show incomer's IP
	p, _ := peer.FromContext(responseStream.Context())
	incomingIP := p.Addr.String()
	log.Printf("Incoming request from IP: %v", incomingIP)

	// declare an AWS client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client = s3.NewFromConfig(cfg)

	// get a file from AWS bucket
	dlInput := &s3.GetObjectInput{
		Bucket: aws.String(req.GetFileBucket()),
		Key:    aws.String(req.GetFileName()),
	}

	file, err := utils.GetFile(context.TODO(), client, dlInput)
	if err != nil {
		log.Printf("GetFile error: %v", err)
		return err
	}

	// save the file
	err = utils.SaveFile(file, req.GetFileName())
	if err != nil {
		log.Printf("SaveFile error: %v\n", err)
		return err
	}

	//send file back by stream
	bufferSize := 64 * 1024 //64KiB, tweak this as desired
	fileUpload, err := os.Open(utils.SplitKeyName(req.GetFileName()))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer fileUpload.Close()
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := fileUpload.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		resp := &proto.Response{
			FileChunk: buff[:bytesRead],
		}
		err = responseStream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
}
