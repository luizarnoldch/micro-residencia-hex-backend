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
	log.Println(request.Body)
	contentType := request.Headers["content-type"]
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Println("Error while getting header: multipart/form-data Error")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Println("Error parsing media type:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var fileData []byte
	if request.IsBase64Encoded {
		fileData, err = base64.StdEncoding.DecodeString(request.Body)
		log.Println("Request is Base64Encoded")
		if err != nil {
			log.Println("Error decoding base64:", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	} else {
		log.Println("Request is not Base64Encoded")
		fileData = []byte(request.Body)
	}

	var fileName string

	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("MediaType: multipart/")
		mr := multipart.NewReader(bytes.NewReader(fileData), params["boundary"])
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Error reading multipart section:", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			if part.FormName() == "file_name" {
				fileName = part.FileName()
				if fileName == "" {
					nameData, err := io.ReadAll(part)
					if err != nil {
						log.Println("Error reading the file_name part:", err)
						return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
					}
					fileName = string(nameData)
				}
			}			
			if part.FileName() != "" {
				log.Println("Filename != ''")
				fileData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading the part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}

				s3client := s3.NewFromConfig(cfg)

				// Aquí obtienes la extensión del archivo
				fileExt := filepath.Ext(part.FileName())

				log.Println("filename: ", fileName)

				// Y aquí la añades a la clave S3
				key := BUCKET_KEY + fileName + fileExt

				input := &s3.PutObjectInput{
					Bucket: aws.String(BUCKET_NAME),
					Key:    aws.String(key),
					Body:   bytes.NewReader(fileData),
				}

				output, err := s3client.PutObject(ctx, input)
				if err != nil {
					log.Println("Error while putting object in S3:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}

				log.Println(output)
			}
		}
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       "file uploaded successfully",
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
