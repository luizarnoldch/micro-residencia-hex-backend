package main

import (
	"context"
	"fmt"
	"os"

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
