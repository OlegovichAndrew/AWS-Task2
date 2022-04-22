package transport

import (
	"aws-server/config"
	"aws-server/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"os"
)

type Server struct {
	proto.UnimplementedAWSServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) DownloadEndpoint(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Downloading started...\n")
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(config.GRPC_ADDR, grpc.WithInsecure())
	if err != nil {
		log.Fatal("client could connect to grpc service:", err)
	}
	log.Printf("gRPC client connected: %v", config.GRPC_ADDR)

	c := proto.NewAWSServiceClient(conn)

	fileStreamResponse, err := c.Download(context.TODO(), &proto.Request{
		FileName:   config.FILE_NAME,
		FileBucket: config.BUCKET_NAME,
	})

	if err != nil {
		log.Println("error downloading:", err)
		return
	}

	f, err := os.Create("filename.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		chunkResponse, err := fileStreamResponse.Recv()
		if err == io.EOF {
			log.Println("received all chunks")
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			break
		}
		if err != nil {
			log.Println("err receiving chunk:", err)
			break
		}
		_, err = f.Write(chunkResponse.FileChunk)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}

		log.Printf("got new chunk with data: %s \n", chunkResponse.FileChunk)
	}
	_, err = fmt.Fprintf(w, "Downloading finished!\n")
	if err != nil {
		log.Println(err)
	}
}
