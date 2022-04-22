package transport

import (
	"aws-http-server/config"
	"aws-http-server/proto"
	"aws-http-server/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	//create connection to grpc server
	conn, err := grpc.Dial(config.GRPC_ADDR, grpc.WithInsecure())
	if err != nil {
		log.Fatal("client could connect to grpc service:", err)
	}
	log.Printf("gRPC client connected: %v", config.GRPC_ADDR)

	c := proto.NewAWSServiceClient(conn)
	//create stream request with env parameters.
	fileStreamResponse, err := c.Download(context.TODO(), &proto.Request{
		FileName:   config.FILE_NAME,
		FileBucket: config.BUCKET_NAME,
	})

	if err != nil {
		log.Println("error downloading:", err)
		return
	}
	//create a new file in which we will write.
	f, err := os.Create(config.FILE_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		chunkResponse, err := fileStreamResponse.Recv()
		if err == io.EOF {
			log.Println("received all chunks")
			//err = f.Close()
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
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
	//increase file value if it's a text file
	if utils.IsTextFile(config.FILE_NAME) {
		err = utils.IncreaseFileValue(config.FILE_NAME)
		if err != nil {
			log.Println(err)
		}
	}

	//send it back to AWS S3 bucket

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	ulFile, err := os.Open(config.FILE_NAME)
	if err != nil {
		fmt.Printf("Unable to open file %v\n", config.FILE_NAME)
		return
	}

	ulInput := &s3.PutObjectInput{
		Bucket: &config.BUCKET_NAME,
		Key:    &config.FILE_NAME,
		Body:   ulFile,
	}

	_, err = client.PutObject(context.TODO(), ulInput)
	if err != nil {
		log.Printf("Got error uploading file:%v\n", err)
		return
	}

	fmt.Fprintf(w, "File sent back to the S3 bucket.")
}
