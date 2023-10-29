package application

import (
	"context"
	"fmt"
	"main/src/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DocumentoServiceDynamo struct {
	client *dynamodb.Client
	table  string
	ctx    context.Context
}

func (dynamo DocumentoServiceDynamo) CreateDocument(req domain.DocumentoRequest) (domain.DocumentoResponse, error) {
	reqToDoc := req.ToDocumento()
	item, err := attributevalue.MarshalMap(reqToDoc)
	if err != nil {
		return domain.DocumentoResponse{Message: err.Error()}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(dynamo.table),
		Item:      item,
	}

	_, err = dynamo.client.PutItem(dynamo.ctx, input)
	if err != nil {
		return domain.DocumentoResponse{Message: err.Error()}, err
	}

	response := reqToDoc.ToDocumentoResponse()

	return response, nil
}

func (dynamo DocumentoServiceDynamo) GetAllDocuments() ([]domain.DocumentoResponse, error) {
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fmt.Sprintf("SELECT * FROM \"%v\"", dynamo.table)),
	}

	response, err := dynamo.client.ExecuteStatement(dynamo.ctx, input)
	if err != nil {
		return nil, err
	}

	var documentos []domain.Documento
	err = attributevalue.UnmarshalListOfMaps(response.Items, &documentos)
	if err != nil {
		return nil, err
	}

	var documentosResponse []domain.DocumentoResponse

	for _, documento := range documentos {
		documentosResponse = append(documentosResponse, documento.ToDocumentoResponse())
	}

	return documentosResponse, nil
}

func (dynamo DocumentoServiceDynamo) UpdateDocument(req domain.DocumentoRequest, id string) (domain.DocumentoResponse, error) {
	reqToDoc := req.ToDocumento()

	update := expression.
		Set(expression.Name("departamento"), expression.Value(reqToDoc.Departamento)).
		Set(expression.Name("residente"), expression.Value(reqToDoc.Residente)).
		Set(expression.Name("fecha_de_pago"), expression.Value(reqToDoc.FechaDePago)).
		Set(expression.Name("tipo_de_servicio"), expression.Value(reqToDoc.TipoDeServicio))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return domain.DocumentoResponse{Message: err.Error()}, err
	}

	id_dynamo, err := attributevalue.Marshal(id)
	if err != nil {
		panic(err)
	}

	key := map[string]types.AttributeValue{"id_documento": id_dynamo}


	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(dynamo.table),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err = dynamo.client.UpdateItem(dynamo.ctx, input)
	if err != nil {
		return domain.DocumentoResponse{Message: err.Error()}, err
	}

	response := reqToDoc.ToDocumentoResponse()

	return response, nil
}

func (dynamo DocumentoServiceDynamo) DeleteDocument(id string) (domain.DocumentoResponse, error) {
	id_dynamo, err := attributevalue.Marshal(id)
	if err != nil {
		panic(err)
	}

	key := map[string]types.AttributeValue{"id_documento": id_dynamo}


	input := &dynamodb.DeleteItemInput{
		TableName:                 aws.String(dynamo.table),
		Key:                       key,
	}

	_, err = dynamo.client.DeleteItem(dynamo.ctx, input)
	if err != nil {
		return domain.DocumentoResponse{Message: err.Error()}, err
	}

	return domain.DocumentoResponse{Message: fmt.Sprintf("Usuario: %s eliminado",id_dynamo)}, nil
}

func NewDocumentoServiceDynamo(client *dynamodb.Client, table string, ctx context.Context) *DocumentoServiceDynamo {
	return &DocumentoServiceDynamo{
		client: client,
		table:  table,
		ctx:    ctx,
	}
}
