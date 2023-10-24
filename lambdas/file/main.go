package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"

	"log"
	"os"
	"strings"

	// // "path/filepath"
	// "net/http"
	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// // "github.com/grokify/go-awslambda"
	// "github.com/olahol/go-imageupload"
)

type CustomStruct struct {
	Content       string
	FileName      string
	FileExtension string
}

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

	fileData, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3client := s3.NewFromConfig(cfg)

	newUUID := uuid.NewString()

	key := BUCKET_KEY + newUUID

	input := &s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
	}

	output, err := s3client.PutObject(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	log.Println(output)

	// fmt.Println(request.Body)
	// pdfData := request.Body

	// response := events.APIGatewayProxyResponse{}

	// pdfBytes, err := base64.StdEncoding.DecodeString(pdfData)
	// if err != nil {
	// 	log.Printf("Error al decodificar el archivo base64: %v", err)
	// 	return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error al procesar el archivo"}, nil
	// }
	// cfg, err := config.LoadDefaultConfig(ctx)

	// if err != nil {
	// 	log.Fatalf("unable to load SDK config, %v", err)
	// }

	// s3client := s3.NewFromConfig(cfg)

	// httpRequest, err := createHTTPRequest(request)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{StatusCode: 500}, err
	// }

	// // r, err := awslambda.NewReaderMultipart(request)
	// // if err != nil {
	// // 	return response, err
	// // }
	// // part, err := r.NextPart()
	// // if err != nil {
	// // 	return response, err
	// // }

	// // content, err := io.ReadAll(part)
	// // if err != nil {
	// // 	return response, err
	// // }

	// img, err := imageupload.Process(httpRequest,"file")

	// if err != nil {
	// 	return response, err
	// }

	// Generar un número aleatorio entre 1 y 100
	// numeroAleatorio := rand.Intn(100) + 1

	// key := BUCKET_KEY + string(numeroAleatorio)

	// input := &s3.PutObjectInput{
	// 	Bucket: aws.String(BUCKET_NAME),
	// 	Key:    aws.String(key),
	// 	Body:   bytes.NewReader(pdfBytes),
	// }

	// output, err := s3client.PutObject(ctx, input)
	// if err != nil {
	// 	return response, err
	// }

	// log.Println(output)

	// custom := CustomStruct{
	// 	Content:       string(content),
	// 	FileName:      part.FileName(),
	// 	FileExtension: filepath.Ext(part.FileName()),
	// }
	// customBytes, err := json.Marshal(custom)
	// if err != nil {
	// 	return response, err
	// }

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       "file upload"}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
