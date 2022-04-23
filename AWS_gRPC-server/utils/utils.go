package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// SplitKeyName - splits the given AWS file KEY and returns just only the file name.
func SplitKeyName(key string) string {
	splitName := strings.Split(key, "/")
	fileName := splitName[len(splitName)-1]
	return fileName
}

// GetFile - downloads the file from AWS.
func GetFile(ctx context.Context, client *s3.Client, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	fmt.Printf("Downloading: %v\n", aws.ToString(input.Key))
	return client.GetObject(ctx, input)
}

func SaveFile(file *s3.GetObjectOutput, key string) error {
	fmt.Printf("Saving: %v\n", key)
	f, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return err
	}

	fileName := SplitKeyName(key)
	return ioutil.WriteFile(fileName, f, 0644)
}
