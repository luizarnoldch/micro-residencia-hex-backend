package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	
	log.Println("Hello world with cognito")

	return "Hello World, Protected", nil
}

func main() {
	lambda.Start(handler)
}
