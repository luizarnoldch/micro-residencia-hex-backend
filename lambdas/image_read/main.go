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
	log.Println("Handler invoked")  // <-- Log at the start

	objectKey := request.QueryStringParameters["key"]
	if objectKey == "" {
		log.Println("Object key not provided")  // <-- Log for missing key
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "Object key is required."}, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	s3client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: &BUCKET_NAME,
		Key:    &objectKey,
	}

	object, err := s3client.GetObject(ctx, input)
	if err != nil {
		log.Println("Error retrieving object:", err)  // <-- Log for error while retrieving object
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	defer object.Body.Close()

	// Read the object's content into a byte slice
	data, err := io.ReadAll(object.Body)
	if err != nil {
		log.Println("Error reading object data:", err)  // <-- Log for error while reading object data
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Convert the byte slice to a base64 encoded string
	encodedString := base64.StdEncoding.EncodeToString(data)

	contentType := "application/octet-stream"
	if object.ContentType != nil {
		log.Println("ContentType != nil")
		contentType = *object.ContentType
	}

	contentDisposition := "attachment"
	if strings.HasSuffix(objectKey, ".jpg") || strings.HasSuffix(objectKey, ".jpeg") {
		log.Println("String has suffix .jpg or .jpeg")
		contentDisposition = "inline"
	} else if strings.HasSuffix(objectKey, ".pdf") {
		log.Println("String has suffix .pdf")
		contentDisposition = "inline; filename=" + objectKey
	}

	headers := map[string]string{
		"Content-Type":                 contentType,
		"Content-Disposition":          contentDisposition,
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	}

	log.Println("Returning response")  // <-- Log before returning the response

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		IsBase64Encoded: true,
		Body:            encodedString, // set the base64 encoded string as the body
	}, nil
}

func main() {
	log.Println("Lambda starting")  // <-- Log at the start of Lambda execution
	lambda.Start(handler)
}
