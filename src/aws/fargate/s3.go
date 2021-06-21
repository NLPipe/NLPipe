package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"path"
)

func uploadFile(uuid string, body []byte) (*s3manager.UploadOutput, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:9090"),
		S3ForcePathStyle: aws.Bool(true),
	})

	file := bytes.NewReader(body)
	svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(svc)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("nlpipe"),
		Key:    aws.String(path.Join("uploads", uuid)),
		Body:   file,
	})

	return result, err
}
