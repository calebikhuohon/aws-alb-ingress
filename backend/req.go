package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

type Result struct {
	Greeting string `json:"greeting"`
}

func GetResult(ctx context.Context, region, bucket, item string) (Result, error) {
	var result Result

	f, err := os.Create(item)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(f,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		return Result{}, err
	}

	b, err := io.ReadAll(f)

	_ = json.Unmarshal(b, &result)
	return result, nil
}
