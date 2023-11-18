package domain

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

var (
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	BUCKET_KEY  = os.Getenv("BUCKET_KEY")
)

type DocumentoRequest struct {
	Departamento   string `json:"departamento"`
	Residente      string `json:"residente"`
	FechaDePago    string `json:"fecha_de_pago"`
	TipoDeServicio string `json:"tipo_de_servicio"`
	StateDocument	string `json:"estado_documento"`
}

type Documento struct {
	Documento_ID   string `dynamodbav:"id_documento" json:"id_documento"`
	Departamento   string `dynamodbav:"departamento" json:"departamento"`
	Residente      string `dynamodbav:"residente" json:"residente"`
	FechaDePago    string `dynamodbav:"fecha_de_pago" json:"fecha_de_pago"`
	TipoDeServicio string `dynamodbav:"tipo_de_servicio" json:"tipo_de_servicio"`
	StateDocument	string `dynamodbav:"estado_documento" json:"estado_documento"`
	UrlPDF         string `dynamodbav:"url_pdf" json:"url_pdf"`
}

func (doc Documento) ToDocumentoResponse() DocumentoResponse {
	return DocumentoResponse{
		Documento_ID:   doc.Documento_ID,
		Departamento:   doc.Departamento,
		Residente:      doc.Residente,
		FechaDePago:    doc.FechaDePago,
		TipoDeServicio: doc.TipoDeServicio,
		StateDocument:  doc.StateDocument,
		UrlPDF:         doc.UrlPDF,
	}
}

func (req DocumentoRequest) ToDocumento() Documento {
	id := uuid.NewString()
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s%s.pdf", BUCKET_NAME, BUCKET_KEY, id)

	return Documento{
		Documento_ID:   id,
		Departamento:   req.Departamento,
		Residente:      req.Residente,
		FechaDePago:    req.FechaDePago,
		TipoDeServicio: req.TipoDeServicio,
		StateDocument:  req.StateDocument,
		UrlPDF:         url,
	}
}

type DocumentoResponse struct {
	Documento_ID   string `json:"id_documento"`
	Departamento   string `json:"departamento"`
	Residente      string `json:"residente"`
	FechaDePago    string `json:"fecha_de_pago"`
	TipoDeServicio string `json:"tipo_de_servicio"`
	UrlPDF         string `json:"url_pdf"`
	StateDocument	string `json:"estado_documento"`
	Message        string `json:"message"`
}
