package application

import (
	"main/src/domain"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DocumentoServiceDynamo struct {
	client *dynamodb.Client
	table  string
}

func (dynamo DocumentoServiceDynamo) CreateDocument(req domain.DocumentoRequest) (domain.DocumentoResponse, error) {
	_, err := attributevalue.MarshalMap(req)
	if err != nil {
		return domain.DocumentoResponse{
			Status:  504,
			Message: err.Error(),
		}, err
	}

	return domain.DocumentoResponse{
		Status:  200,
		Message: "item guardado",
	}, nil
}


func NewDocumentoServiceDynamo(client *dynamodb.Client, table  string) *DocumentoServiceDynamo {
	return &DocumentoServiceDynamo{
		client: client,
		table: table,
	}
}