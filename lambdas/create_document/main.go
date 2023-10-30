package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"log"

	"main/src/application"
	"main/src/domain"
	"main/src/infrastructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	TABLE_NAME = os.Getenv("TABLE_NAME")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Inicio de la función Lambda")

	log.Println("Undecoded: ", request.Body)

	dynamoClient, err := infrastructure.GetDynamoClient(ctx)
	if err != nil {
		log.Println("Failed to get dynamodb client:", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to get dynamodb client %s", err),
			StatusCode: 504}, nil
	}

	var documentoRequest domain.DocumentoRequest

	log.Println("Decodificando el cuerpo de la solicitud...")

	log.Println("Undecoded: ", request.Body)

	decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Println("Error decoding base64 request body:", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Error decoding base64: %s", err), StatusCode: 400}, nil
	}

	log.Println("Decoded: ", decodedBody)

	log.Println("Convirtiendo el cuerpo decodificado a JSON...")
	if err := json.Unmarshal([]byte(request.Body), &documentoRequest); err != nil {
		log.Println("Error parsing request body as JSON:", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 502}, nil
	}

	log.Println("Creando documento en la base de datos...")
	dynamoService := application.NewDocumentoServiceDynamo(dynamoClient, TABLE_NAME, ctx)
	response, err := dynamoService.CreateDocument(documentoRequest)
	if err != nil {
		log.Printf("Error creating documento in database: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 400}, nil
	}

	log.Println("Convirtiendo la respuesta a JSON...")
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response to JSON: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 500}, nil
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST,OPTIONS,DELETE,GET,HEAD,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Custom-Header",
		"Content-Type":                 "application/json",
	}

	log.Println("Finalizando la función Lambda con éxito")
	return events.APIGatewayProxyResponse{
		Headers: headers,
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
