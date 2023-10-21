package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grokify/go-awslambda"
)

type CustomStruct struct {
	Content       string
	FileName      string
	FileExtension string
}

var (
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	BUCKET_KEY = os.Getenv("BUCKET_KEY")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3client := s3.NewFromConfig(cfg)

	response := events.APIGatewayProxyResponse{}
	r, err := awslambda.NewReaderMultipart(request)
	if err != nil {
		return response, err
	}
	part, err := r.NextPart()
	if err != nil {
		return response, err
	}

	input := s3.PutObjectInput{
        Bucket: aws.String(BUCKET_NAME),
        Key:    aws.String(BUCKET_KEY),
        Body:   part,
    }

	output, err := s3client.PutObject(ctx, &input)
	if err != nil {
		return response, err
	}

	log.Println(output)
	

	content, err := io.ReadAll(part)
	if err != nil {
		return response, err
	}
	custom := CustomStruct{
		Content:       string(content),
		FileName:      part.FileName(),
		FileExtension: filepath.Ext(part.FileName()),
	}
	customBytes, err := json.Marshal(custom)
	if err != nil {
		return response, err
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	response = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(customBytes)}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
