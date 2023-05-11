package main

type EnvConfig struct {
	Ports struct {
		HTTP string `envconfig:"HTTP_PORT" default:":8080"`
	}
	AWSRegion string `envconfig:"AWS_REGION" default:"us-east-2"`
	S3Bucket  string `envconfig:"S3_BUCKET" default:"ikh-json-bucket"`
	S3Item    string `envconfig:"S3_ITEM" default:"file.json"`
}

func (c EnvConfig) Validate() error {
	return nil
}
