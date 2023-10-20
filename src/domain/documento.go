package domain

import "github.com/google/uuid"

type DocumentoRequest struct {
	Departamento   string `json:"departamento"`
	Residente      string `json:"residente"`
	FechaDePago    string `json:"fecha_de_pago"`
	TipoDeServicio string `json:"tipo_de_servicio"`
}

type Documento struct {
	Documento_ID   string `dynamodbav:"document_id"`
	Departamento   string `dynamodbav:"departamento"`
	Residente      string `dynamodbav:"residente"`
	FechaDePago    string `dynamodbav:"fecha_de_pago"`
	TipoDeServicio string `dynamodbav:"tipo_de_servicio"`
}

func (req DocumentoRequest) ToDocumento() Documento {
	return Documento{
		Documento_ID:   uuid.NewString(),
		Departamento:   req.Departamento,
		Residente:      req.Residente,
		FechaDePago:    req.FechaDePago,
		TipoDeServicio: req.TipoDeServicio,
	}
}

type DocumentoSimpleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type DocumentoResponse struct {
	Documento_ID   string `json:"id_documento"`
	Departamento   string `json:"departamento"`
	Residente      string `json:"residente"`
	FechaDePago    string `json:"fecha_de_pago"`
	TipoDeServicio string `json:"tipo_de_servicio"`
}
