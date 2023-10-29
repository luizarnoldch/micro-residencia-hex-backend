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
	dynamoClient, err := infrastructure.GetDynamoClient(ctx)
	if err != nil {
		log.Fatalln("Failed to get dynamodb client")
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to get dynamodb client %s", err),
			StatusCode: 504}, nil
	}

	var documentoRequest domain.DocumentoRequest

	decodedBody, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Println("Error decoding base64 request body.")
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Error decoding base64: %s", err), StatusCode: 400}, nil
	}

	if err := json.Unmarshal(decodedBody, &documentoRequest); err != nil {
		log.Println("Error parsing request body as JSON.")
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 502}, nil
	}

	dynamoService := application.NewDocumentoServiceDynamo(dynamoClient, TABLE_NAME, ctx)

	id_documento := request.PathParameters["id_documento"]

	response, err := dynamoService.UpdateDocument(documentoRequest,id_documento)
	if err != nil {
		log.Printf("error creating documento in database :%s\n", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 400}, nil
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("error marshaling response to JSON: %s\n", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 500}, nil
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type": "application/json",
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
