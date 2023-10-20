package application

import (
	"context"
	"fmt"
	"main/src/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DocumentoServiceDynamo struct {
	client *dynamodb.Client
	table  string
	ctx    context.Context
}

func (dynamo DocumentoServiceDynamo) CreateDocument(req domain.DocumentoRequest) (domain.DocumentoSimpleResponse, error) {
    reqToDoc := req.ToDocumento()
    item, err := attributevalue.MarshalMap(reqToDoc)
    if err != nil {
        return domain.DocumentoSimpleResponse{
            Status:  504,
            Message: err.Error(),
        }, err
    }

    input := &dynamodb.PutItemInput{
        TableName: aws.String(dynamo.table),
        Item:      item,
    }

    _, err = dynamo.client.PutItem(dynamo.ctx, input)
    if err != nil {
        return domain.DocumentoSimpleResponse{
            Status:  503,
            Message: err.Error(),
        }, err
    }

    return domain.DocumentoSimpleResponse{
        Status:  200,
        Message: "item guardado",
    }, nil
}

func (dynamo DocumentoServiceDynamo) GetAllDocuments() ([]domain.DocumentoResponse, error) {
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fmt.Sprintf("SELECT * FROM \"%v\"", dynamo.table)),
	}

	response, err := dynamo.client.ExecuteStatement(dynamo.ctx, input)
	if err != nil {
		return nil, err
	}

	var documentos []domain.DocumentoResponse
	err = attributevalue.UnmarshalListOfMaps(response.Items, &documentos)
	if err != nil {
		return nil, err
	}

	return documentos, nil
}

func NewDocumentoServiceDynamo(client *dynamodb.Client, table string, ctx context.Context) *DocumentoServiceDynamo {
	return &DocumentoServiceDynamo{
		client: client,
		table:  table,
		ctx:    ctx,
	}
}
