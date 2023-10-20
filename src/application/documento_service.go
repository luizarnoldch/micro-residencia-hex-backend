package application

import "main/src/domain"

type DocumentoService interface {
	CreateDocument(domain.DocumentoRequest) (domain.DocumentoResponse, error)
	GetAllDocuments() ([]domain.DocumentoResponse, error)
}