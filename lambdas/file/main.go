package main

import (
	"context"
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/grokify/go-awslambda"
)

type CustomStruct struct {
	Content       string
	FileName      string
	FileExtension string
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	response := events.APIGatewayProxyResponse{}
	r, err := awslambda.NewReaderMultipart(request)
	if err != nil {
		return response, err
	}
	part, err := r.NextPart()
	if err != nil {
		return response, err
	}
	content, err := io.ReadAll(part)
	if err != nil {
		return response, err
	}
	custom := CustomStruct{
		Content:       string(content),
		FileName:      part.FileName(),
		FileExtension: filepath.Ext(part.FileName())}

	customBytes, err := json.Marshal(custom)
	if err != nil {
		return response, err
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	response = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(customBytes)}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
