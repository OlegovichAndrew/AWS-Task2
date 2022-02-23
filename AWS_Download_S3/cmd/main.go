package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"aws-dl-s3/proto"
	"aws-dl-s3/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	bucket, fileKey string
	client          *s3.Client
	addr            = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.StringVar(&bucket, "b", "", "The bucket to download/upload the file from/to")
	flag.StringVar(&fileKey, "f", "", "The file to download/upload")
	flag.Parse()

	if bucket == "" || fileKey == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
		return
	}

	// declare a gRPC client
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	gRPCClient := proto.NewAWSServiceClient(conn)

	// declare an AWS client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client = s3.NewFromConfig(cfg)

	// get a file from AWS bucket
	dlInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	file, err := utils.GetFile(context.TODO(), client, dlInput)
	if err != nil {
		log.Printf("GetFile error: %v", err)
		return
	}

	// save the file
	err = utils.SaveFile(file, fileKey)
	if err != nil {
		log.Printf("SaveFile error: %v\n", err)
		return
	}

	//send a message by gRPC
	r, err := gRPCClient.StringSend(context.TODO(), &proto.Request{
		TestReq: fmt.Sprintf("File %v was saved successfuly", utils.SplitKeyName(fileKey)),
	})

	if err != nil {
		log.Printf("couldn't send gRPC request: %v\n", err)
	}
	log.Printf("Greeting: %s\n", r.GetTestResp())
}
