package application

import "main/src/domain"

type DocumentoService interface {
	CreateDocument(domain.DocumentoRequest) (domain.DocumentoSimpleResponse, error)
	GetAllDocuments() ([]domain.DocumentoResponse, error)
	UpdateDocument(domain.DocumentoRequest,string) (domain.DocumentoSimpleResponse, error)
	DeleteDocument(string) (domain.DocumentoSimpleResponse, error)
}