package main

import (
	"context"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"log"

	"main/src/application"
	"main/src/domain"
	"main/src/infrastructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	TABLE_NAME = os.Getenv("TABLE_NAME")
	SQS_NAME = os.Getenv("SQS_NAME")
)

// Definir una estructura que represente el objeto JSON
type FileData struct {
	FileContents string `json:"file_contents"` // Contenido del archivo codificado en base64
	FileName     string `json:"file_name"`     // Nombre real del archivo
	RealFileName     string `json:"real_file_name"`     // Nombre real del archivo
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Inicio de la función Lambda")

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 504,
			Body:       err.Error(),
		}, fmt.Errorf("error loading AWS configuration: %w", err)
	}

	dynamoClient, err := infrastructure.GetDynamoClient(ctx)
	if err != nil {
		log.Println("Failed to get dynamodb client:", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to get dynamodb client %s", err),
			StatusCode: 504}, nil
	}

	log.Println("Starting Lambda handler")
	contentType := request.Headers["content-type"]
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Println("Error: content type not multipart/form-data")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Println("Error parsing media type:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	
	log.Println("Request.Body:")
	log.Println(request.Body)

	var fileData []byte
	if request.IsBase64Encoded {
		log.Println("Request body is Base64Encoded")
		fileData, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Println("Error decoding base64 body:", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	} else {
		log.Println("Request body is not Base64Encoded")
		fileData = []byte(request.Body)
	}

	log.Println("fileData Based64Encoded:")
	log.Println(fileData)

	var fileName string
	var fileDepartamento string
	var fileResidente string
	var fileFechaPago string
	var fileTipoServicio string
	var fileStateDocument string 
	var fileBuffer bytes.Buffer
	var realFileName string 
	

	if strings.HasPrefix(mediaType, "multipart/") {
		// Crea un multipart reader
		reader := multipart.NewReader(bytes.NewReader(fileData), params["boundary"])
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				log.Println("Reached end of multipart content")
				break
			}
			if err != nil {
				log.Println("Error reading multipart section:", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
			switch part.FormName() {
			case "file_name":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading file_name part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileName = string(nameData)
				log.Println("Received file name:", fileName)
			case "file":
				log.Println("Reading file content")
				realFileName = part.FileName()
				if _, err := io.Copy(&fileBuffer, part); err != nil {
					log.Println("Error copying file content to buffer:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
			case "departamento":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading departamento part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileDepartamento = string(nameData)
				log.Println("Received departamento name:", fileDepartamento)
			case "residente":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading residente part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileResidente = string(nameData)
				log.Println("Received residente name:", fileResidente)
			case "fecha_de_pago":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading fecha_de_pago part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileFechaPago = string(nameData)
				log.Println("Received fecha_de_pago:", fileFechaPago)
			case "tipo_de_servicio":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading tipo_de_servicio part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileTipoServicio = string(nameData)
				log.Println("Received tipo_de_servicio:", fileTipoServicio)
			
			case "estado_documento":
				nameData, err := io.ReadAll(part)
				if err != nil {
					log.Println("Error reading tipo_de_servicio part:", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				fileStateDocument = string(nameData)
				log.Println("Received tipo_de_servicio:", fileStateDocument)
			}
		}
	}

	log.Println("fileName: ",fileName)
	log.Println("fileDepartamento: ",fileDepartamento)
	log.Println("fileResidente: ",fileResidente)
	log.Println("fileFechaPago: ",fileFechaPago)
	log.Println("fileTipoServicio: ",fileTipoServicio)
	log.Println("fileStateDocument: ",fileStateDocument)
	log.Println("fileBuffer: ",fileBuffer)
	log.Println("realFileName: ",realFileName)
	
	

	documentoRequest := domain.DocumentoRequest{
		Departamento: fileDepartamento,
		Residente: fileResidente,
		FechaDePago: fileFechaPago,
		TipoDeServicio: fileTipoServicio,
		StateDocument: fileStateDocument,
	}

	// Convertir el bytes.Buffer a base64 para que pueda ser representado en JSON
	encodedFileContents := base64.StdEncoding.EncodeToString(fileBuffer.Bytes())

	log.Println("Creando documento en la base de datos...")
	dynamoService := application.NewDocumentoServiceDynamo(dynamoClient, TABLE_NAME, ctx)
	response, err := dynamoService.CreateDocument(documentoRequest)
	if err != nil {
		log.Printf("Error creating documento in database: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 400}, nil
	}

	// Crear una instancia de la estructura con los datos codificados y el nombre del archivo
	newData := FileData{
		FileContents: encodedFileContents,
		FileName:     realFileName,
		RealFileName: response.Documento_ID,
	}

	// Serializar la estructura a JSON
	jsonData, err := json.Marshal(newData)
	if err != nil {
		fmt.Println("Error al serializar a JSON:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	// Convertir los datos JSON a una cadena y mostrarla
	jsonString := string(jsonData)
	fmt.Println(jsonString)

	// ======================== SQS code ==========================================

	// Create an Amazon SQS client
	sqsClient := sqs.NewFromConfig(cfg)

	// Send message
	msgInput := &sqs.SendMessageInput{
		QueueUrl:               aws.String(SQS_NAME),
		MessageBody:            aws.String(jsonString),
	}
	
	_, err = sqsClient.SendMessage(ctx, msgInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 504,
			Body:       err.Error(),
		}, fmt.Errorf("error sending SQS message: %w", err)
	}

	log.Println("Convirtiendo la respuesta a JSON...")
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response to JSON: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", err), StatusCode: 500}, nil
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST,OPTIONS,DELETE,GET,HEAD,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Custom-Header",
		"Content-Type":                 "application/json",
	}

	log.Println("Finalizando la función Lambda con éxito")
	return events.APIGatewayProxyResponse{
		Headers: headers,
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
