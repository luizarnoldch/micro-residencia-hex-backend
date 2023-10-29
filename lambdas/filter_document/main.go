package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/src/domain"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

	departamento := request.QueryStringParameters["departamento"]
	residente := request.QueryStringParameters["residente"]
	fechaDePago := request.QueryStringParameters["fecha_de_pago"]

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

	var documentos []domain.Documento
	err = attributevalue.UnmarshalListOfMaps(output.Items, &documentos)
	if err != nil {
		return errorResponse(fmt.Sprintf("Failed parsing DynamoDB response: %s", err)), nil
	}

	var documentosResponse []domain.DocumentoResponse

	for _, documento := range documentos {
		documentosResponse = append(documentosResponse, documento.ToDocumentoResponse())
	}

	body, err := json.Marshal(documentosResponse)
	if err != nil {
		log.Printf("Failed to marshal response: %s", err)
		return errorResponse(fmt.Sprintf("Failed to marshal response: %s", err)), nil
	}

	log.Println("Handler completed successfully")
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
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
