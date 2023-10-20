package main

import (
	"context"
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
	dynamoClient, err := infrastructure.GetDynamoClient(ctx)
	if err != nil {
		log.Fatalln("Failed to get dynamodb client")
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to get dynamodb client %s", err),
			StatusCode: 500}, nil
	}

	var documentoRequest domain.DocumentoRequest

	if err := json.Unmarshal([]byte(request.Body), &documentoRequest); err != nil {
		log.Println("Error parsing request body as JSON.")
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err)}, nil
	}

	dynamoService := application.NewDocumentoServiceDynamo(dynamoClient, TABLE_NAME)

	response, err := dynamoService.CreateDocument(documentoRequest)
	if err != nil {
		log.Printf("error creating documento in database :%s\n", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 400}, nil
	}

	log.Println(response)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "Hello World",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
