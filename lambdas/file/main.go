package lambdas

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
	res := events.APIGatewayProxyResponse{}
	r, err := awslambda.NewReaderMultipart(request)
	if err != nil {
		return res, err
	}
	part, err := r.NextPart()
	if err != nil {
		return res, err
	}
	content, err := io.ReadAll(part)
	if err != nil {
		return res, err
	}
	custom := CustomStruct{
		Content:       string(content),
		FileName:      part.FileName(),
		FileExtension: filepath.Ext(part.FileName())}

	customBytes, err := json.Marshal(custom)
	if err != nil {
		return res, err
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Content-Type":                "application/json",
	}

	res = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(customBytes)}
	return res, nil
}

func main() {
	lambda.Start(handler)
}
