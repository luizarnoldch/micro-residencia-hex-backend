package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	TABLE_NAME = os.Getenv("TABLE_NAME")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Starting the handler")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load SDK config: %s", err)
		return errorResponse(fmt.Sprintf("Failed to load SDK config: %s", err)), nil
	}

	client := dynamodb.NewFromConfig(cfg)

	decodedInput, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Printf("Failed to decode base64 input: %s", err)
		return errorResponse(fmt.Sprintf("Failed to decode base64 input: %s", err)), nil
	}

	// Parse the JSON input
	var inputData struct {
		Departamento string `json:"departamento"`
		Residente    string `json:"residente"`
		FechaDePago  string `json:"fecha_de_pago"`
	}
	if err := json.Unmarshal(decodedInput, &inputData); err != nil {
		log.Printf("Failed to parse JSON input: %s", err)
		return errorResponse(fmt.Sprintf("Failed to parse JSON input: %s", err)), nil
	}
	
	departamento := inputData.Departamento
	residente := inputData.Residente
	fechaDePago := inputData.FechaDePago

	log.Printf("Received filters: departamento: %s, residente: %s, fechaDePago: %s", departamento, residente, fechaDePago)

	filterExpression := ""
	expressionAttributeValues := map[string]types.AttributeValue{}

	if departamento != "" {
		filterExpression = "departamento = :departamentoVal"
		expressionAttributeValues[":departamentoVal"] = &types.AttributeValueMemberS{Value: departamento}
	}

	if residente != "" {
		if filterExpression != "" {
			filterExpression += " AND "
		}
		filterExpression += "residente = :residenteVal"
		expressionAttributeValues[":residenteVal"] = &types.AttributeValueMemberS{Value: residente}
	}

	if fechaDePago != "" {
		if filterExpression != "" {
			filterExpression += " AND "
		}
		filterExpression += "fecha_de_pago = :fechaDePagoVal"
		expressionAttributeValues[":fechaDePagoVal"] = &types.AttributeValueMemberS{Value: fechaDePago}
	}

	log.Printf("Constructed filter expression: %s", filterExpression)

	queryInput := &dynamodb.ScanInput{
		TableName:                 &TABLE_NAME,
		FilterExpression:          &filterExpression,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	output, err := client.Scan(ctx, queryInput)
	if err != nil {
		log.Printf("Failed to scan DynamoDB: %s", err)
		return errorResponse(fmt.Sprintf("Failed to scan DynamoDB: %s", err)), nil
	}

	body, err := json.Marshal(output.Items)
	if err != nil {
		log.Printf("Failed to marshal response: %s", err)
		return errorResponse(fmt.Sprintf("Failed to marshal response: %s", err)), nil
	}

	log.Println("Handler completed successfully")
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func errorResponse(err string) events.APIGatewayProxyResponse {
	log.Printf("Returning error response: %s", err)
	return events.APIGatewayProxyResponse{
		Body:       err,
		StatusCode: 500,
	}
}

func main() {
	lambda.Start(handler)
}
