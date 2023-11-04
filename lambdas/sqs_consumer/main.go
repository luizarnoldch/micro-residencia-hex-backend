package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// FileData structure to match the JSON structure
type FileData struct {
	FileContents string `json:"file_contents"` // Contenido del archivo codificado en base64
	FileName     string `json:"file_name"`     // Nombre real del archivo
	RealFileName     string `json:"real_file_name"`     // Nombre real del archivo
}
// FileData structure to match the JSON structure
var (
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	BUCKET_KEY  = os.Getenv("BUCKET_KEY")
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("SQS Lambda start")

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create an S3 service client
	s3client := s3.NewFromConfig(cfg)

	for _, message := range sqsEvent.Records {
		log.Printf("Processing message %s for event source %s\n", message.MessageId, message.EventSource)

		// Parse the JSON body into the FileData structure
		var fileData FileData
		err := json.Unmarshal([]byte(message.Body), &fileData)
		if err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		// Decode the base64 string to []byte
		fileContent, err := base64.StdEncoding.DecodeString(fileData.FileContents)
		if err != nil {
			log.Printf("Error decoding base64: %v\n", err)
			continue
		}

		// Log the file name
		log.Println("File name:", fileData.FileName)
		// Log the file name
		log.Println("Real File name:", fileData.RealFileName)

		// If the file content is not empty, store it in the S3 bucket
		if len(fileContent) > 0 && fileData.FileName != "" {
			log.Println("Storing file to S3 bucket")
			fileExt := filepath.Ext(fileData.RealFileName)
			key := BUCKET_KEY + fileData.FileName + fileExt
			log.Println("Using S3 key:", key)

			// Create a reader from the file content
			reader := bytes.NewReader(fileContent)

			// Prepare the PutObject input
			input := &s3.PutObjectInput{
				Bucket: aws.String(BUCKET_NAME),
				Key:    aws.String(key),
				Body:   reader,
			}

			// Upload the file to S3
			output, err := s3client.PutObject(ctx, input)
			if err != nil {
				log.Printf("Error while putting object to S3: %v\n", err)
				continue
			}
			log.Printf("Successfully stored to S3: %v\n", output)
		} else {
			log.Println("File content is empty or file name is missing")
		}
	}

	log.Println("SQS Lambda end")
	return nil
}

func main() {
	lambda.Start(handler)
}
