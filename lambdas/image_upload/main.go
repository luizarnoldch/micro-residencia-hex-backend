package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	BUCKET_KEY  = os.Getenv("BUCKET_KEY")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Starting Lambda handler")
	contentType := request.Headers["content-type"]
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Println("Error: content type not multipart/form-data")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Println("Error parsing media type:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading SDK config: %v", err)
	}

	var fileData []byte
	if request.IsBase64Encoded {
		log.Println("Request body is Base64Encoded")
		fileData, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Println("Error decoding base64 body:", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	} else {
		log.Println("Request body is not Base64Encoded")
		fileData = []byte(request.Body)
	}

	var fileName string
	var fileBuffer bytes.Buffer

	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("Handling multipart content")
		mr := multipart.NewReader(bytes.NewReader(fileData), params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				log.Println("Reached end of multipart content")
				break
			}
			if err != nil {
				log.Println("Error reading multipart section:", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			switch part.FormName() {
			case "file_name":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading file_name part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileName = string(nameData)
				log.Println("Received file name:", fileName)
			case "file":
				log.Println("Reading file content")
				if _, err := io.Copy(&fileBuffer, part); err != nil {
					log.Println("Error copying file content to buffer:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
			}
		}
	}

	if fileBuffer.Len() > 0 && fileName != "" {
		log.Println("Storing file to S3 bucket")
		s3client := s3.NewFromConfig(cfg)
		fileExt := filepath.Ext(fileName)
		key := BUCKET_KEY + fileName + fileExt
		log.Println("Using S3 key:", key)

		input := &s3.PutObjectInput{
			Bucket: aws.String(BUCKET_NAME),
			Key:    aws.String(key),
			Body:   bytes.NewReader(fileBuffer.Bytes()),
		}

		output, err := s3client.PutObject(ctx, input)
		if err != nil {
			log.Println("Error while putting object to S3:", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		log.Println("Successfully stored to S3:", output)
	} else {
		log.Println("Either file buffer is empty or filename is missing")
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       "file uploaded successfully",
	}, nil
}

func main() {
	log.Println("Starting Lambda execution")
	lambda.Start(handler)
}
