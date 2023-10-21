package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"log"

	"main/src/application"
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

	dynamoService := application.NewDocumentoServiceDynamo(dynamoClient, TABLE_NAME, ctx)

	response, err := dynamoService.GetAllDocuments()
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
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
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
