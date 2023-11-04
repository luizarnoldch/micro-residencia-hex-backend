package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) (error) {
	log.Println("sqs lambda start")
	for _, message := range sqsEvent.Records {
		log.Printf("The message %s for event source %s = %s\n", message.MessageId, message.EventSource, message.Body)
        fmt.Printf("The message %s for event source %s = %s\n", message.MessageId, message.EventSource, message.Body)
    }

	log.Println("sqs lambda end")

	// headers := map[string]string{
	// 	"Access-Control-Allow-Origin": "*",
	// 	"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
	// 	"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	// 	"Content-Type": "application/json",
	// }

	return nil
}

func main() {
	lambda.Start(handler)
}
