package main

import (
	"context"
	"fmt"
	"os"
	// "mime"
	"encoding/base64"
	// "mime/multipart"
	"log"
	// "strings"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	SQS_NAME = os.Getenv("SQS_NAME")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 504,
			Body:       err.Error(),
		}, fmt.Errorf("error loading AWS configuration: %w", err)
	}

	// log.Println("Starting Lambda handler")
	// contentType := request.Headers["content-type"]
	// if !strings.Contains(contentType, "multipart/form-data") {
	// 	log.Println("Error: content type not multipart/form-data")
	// 	return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	// }

	// mediaType, params, err := mime.ParseMediaType(contentType)
	// if err != nil {
	// 	log.Println("Error parsing media type:", err)
	// 	return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	// }
	
	log.Println("Request.Body:")
	log.Println(request.Body)

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

	log.Println("fileData Based64Encoded:")
	log.Println(fileData)

	// var fileName string
	// var fileBuffer bytes.Buffer
	// var realFileName string 

	// if strings.HasPrefix(mediaType, "multipart/") {
	// 	// Crea un multipart reader
	// 	reader := multipart.NewReader(strings.NewReader(request.Body), params["boundary"])

	// 	// Parsea todos los campos
	// 	form, err := reader.ReadForm(32 << 20) // 32MB es el tamaño máximo de la memoria
	// 	if err != nil {
	// 		fmt.Println("Error en ReadForm:", err)
	// 		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	// 	}
	// 	defer form.RemoveAll()

	// 	// Procesar formularios y archivos
	// 	for key, values := range form.Value {
	// 		for _, value := range values {
	// 			fmt.Printf("Key: %v, Value: %v\n", key, value)
	// 		}
	// 	}

	// 	// Aquí podrías procesar los archivos, guardarlos en S3, etc.
	// 	// En este ejemplo, simplemente los imprimimos en la salida estándar (stdout).
	// 	for _, fileHeaders := range form.File {
	// 		for _, fileHeader := range fileHeaders {
	// 			file, err := fileHeader.Open()
	// 			if err != nil {
	// 				fmt.Println("Error al abrir el archivo:", err)
	// 				continue
	// 			}
	// 			defer file.Close()

	// 			fmt.Printf("Procesando archivo: %v\n", fileHeader.Filename)
	// 			// Haz algo con el archivo, como leerlo o subirlo a S3...
	// 			// Este es solo un ejemplo para leer y mostrar los primeros 512 bytes del archivo
	// 			buffer := make([]byte, 512)
	// 			n, err := file.Read(buffer)
	// 			if err != nil && err != io.EOF {
	// 				fmt.Println("Error al leer el archivo:", err)
	// 				continue
	// 			}
	// 			fmt.Printf("Leídos %d bytes del archivo: %s\n", n, string(buffer[:n]))
	// 		}
	// 	}
	// }




	// ======================== SQS code ==========================================

	// Create an Amazon SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Send message
	msgInput := &sqs.SendMessageInput{
		QueueUrl:               aws.String(SQS_NAME),
		MessageBody:            aws.String("Hello World!"),
	}
	

	_, err = sqsClient.SendMessage(ctx, msgInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 504,
			Body:       err.Error(),
		}, fmt.Errorf("error sending SQS message: %w", err)
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
		Body:       "Hello World",
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
