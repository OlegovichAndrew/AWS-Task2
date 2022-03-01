package transport

import (
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
	fmt.Fprintf(w, "Downloading started\n")

	conn, err := grpc.Dial("localhost:4444", grpc.WithInsecure())
	if err != nil {
		log.Fatal("client could connect to grpc service:", err)
	}
	c := proto.NewAWSServiceClient(conn)

	fileStreamResponse, err := c.Download(context.TODO(), &proto.Request{
		FileName:   "The_first_upload/number.txt",
		FileBucket: "upload.practice",
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
}
