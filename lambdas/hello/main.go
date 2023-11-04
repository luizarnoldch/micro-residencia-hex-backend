package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var (
	SQS_NAME = os.Getenv("SQS_NAME")
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 504,
			Body:       err.Error(),
		}, fmt.Errorf("error loading AWS configuration: %w", err)
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
			}
		}
	}

	fmt.Println("fileName: ",fileName)
	fmt.Println("fileDepartamento: ",fileDepartamento)
	fmt.Println("fileResidente: ",fileResidente)
	fmt.Println("fileFechaPago: ",fileFechaPago)
	fmt.Println("fileTipoServicio: ",fileTipoServicio)
	fmt.Println("fileBuffer: ",fileBuffer)
	fmt.Println("realFileName: ",realFileName)


	type FileInformation struct {
		FileName        string `json:"file_name"`
		Departamento    string `json:"departamento"`
		Residente       string `json:"residente"`
		FechaDePago     string `json:"fecha_de_pago"`
		TipoDeServicio  string `json:"tipo_de_servicio"`
		RealFileName    string `json:"real_file_name"`
	}

	// Crear una instancia de la estructura y asignar los valores de las variables
	fileInfo := FileInformation{
		FileName:        fileName,
		Departamento:    fileDepartamento,
		Residente:       fileResidente,
		FechaDePago:     fileFechaPago,
		TipoDeServicio:  fileTipoServicio,
		RealFileName:    realFileName,
	}

	// Serializar la estructura a JSON
	jsonData, err := json.Marshal(fileInfo)
	if err != nil {
		fmt.Println(err)
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

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "DELETE,GET,HEAD,POST,PUT",
		"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Content-Type":                 "application/json",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       "Hello World",
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
