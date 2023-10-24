package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	objectKey := request.QueryStringParameters["key"]
	if objectKey == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Object key is required."}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: &BUCKET_NAME,
		Key:    &objectKey,
	}

	object, err := s3client.GetObject(ctx, input)
	if err != nil {
		log.Println("Error retrieving object:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	defer object.Body.Close()

	// Read the object's content into a byte slice
	data, err := io.ReadAll(object.Body)
	if err != nil {
		log.Println("Error reading object data:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Convert the byte slice to a base64 encoded string
	encodedString := base64.StdEncoding.EncodeToString(data)

	contentType := "application/octet-stream"
	if object.ContentType != nil {
		contentType = *object.ContentType
	}

	contentDisposition := "attachment"
	if strings.HasSuffix(objectKey, ".jpg") || strings.HasSuffix(objectKey, ".jpeg") {
		contentDisposition = "inline"
	}

	headers := map[string]string{
		"Content-Type":                 contentType,
		"Content-Disposition":          contentDisposition,
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		IsBase64Encoded: true,
		Body:            encodedString, // set the base64 encoded string as the body
	}, nil
}

func main() {
	lambda.Start(handler)
}
