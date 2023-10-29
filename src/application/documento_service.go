package application

import "main/src/domain"

type DocumentoService interface {
	CreateDocument(domain.DocumentoRequest) (domain.DocumentoResponse, error)
	GetAllDocuments() ([]domain.DocumentoResponse, error)
	UpdateDocument(domain.DocumentoRequest,string) (domain.DocumentoResponse, error)
	DeleteDocument(string) (domain.DocumentoResponse, error)
}