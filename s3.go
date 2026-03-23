package main

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getS3Client() (s3.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	s3Client := *s3.NewFromConfig(sdkConfig)

	return s3Client, err
}

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions used in the examples. It contains S3Client, an Amazon S3 service client that is used to perform bucket and object actions.
type BucketBasics struct {
	S3Client *s3.Client
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics BucketBasics) UploadFile(body io.Reader, filename string) (string, error) {
	_, err := basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        body,
		ContentType: aws.String("json"),
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

// DownloadFile gets an object from a bucket and returns its contents.
func (basics BucketBasics) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, filename, err)
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", filename, err)
		return nil, err
	}

	return body, nil
}
