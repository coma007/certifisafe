package dto

import (
	"certifisafe-back/model"
	"time"
)

type RequestDTO struct {
	ParentCertificate *CertificateDTO
	Certificate       *CertificateDTO
	Datetime          time.Time
	Status            string
}

type NewRequestDTO struct {
	ParentCertificate *CertificateDTO
	Certificate       *CertificateDTO
	Datetime          time.Time
}

func RequestToDTO(req *model.Request) *RequestDTO {
	request := RequestDTO{
		CertificateToDTO(req.ParentCertificate),
		CertificateToDTO(req.Certificate),
		req.Datetime,
		RequestStatusToString(req.Status),
	}
	return &request
}

func RequestDTOtoModel(request *RequestDTO) *model.Request {
	return &model.Request{
		0,
		CertificateDTOtoModel(request.ParentCertificate),
		CertificateDTOtoModel(request.Certificate),
		request.Datetime,
		StringToRequestStatus(request.Status)}
}

func NewRequestDTOtoModel(request *NewRequestDTO) *model.Request {
	return &model.Request{0,
		CertificateDTOtoModel(request.ParentCertificate),
		CertificateDTOtoModel(request.ParentCertificate),
		request.Datetime,
		model.RequestStatus(model.PENDING)}
}

func RequestStatusToString(reqStatus model.RequestStatus) string {
	switch reqStatus {
	case model.PENDING:
		return "PENDING"
	case model.ACCEPTED:
		return "ACCEPTED"
	case model.REJECTED:
		return "REJECTED"
	}
	// TODO error
	return ""
}

func StringToRequestStatus(certStatus string) model.RequestStatus {
	switch certStatus {
	case "PENDING":
		return model.PENDING
	case "ACCEPTED":
		return model.ACCEPTED
	case "REJECTED":
		return model.REJECTED
	}
	// TODO error
	return -1
}
