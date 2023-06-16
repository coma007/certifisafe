package request

import (
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"time"
)

type RequestDTO struct {
	ParentCertificate *certificate.CertificateDTO
	CertificateName   string
	CertificateType   string
	Subject           *user.UserBaseDTO
	Datetime          time.Time
	Status            string
}

type NewRequestDTO struct {
	ParentSerial    *uint
	CertificateName string
	CertificateType string
	Token           string
}

func RequestToDTO(req *Request) *RequestDTO {
	if req == nil {
		return &RequestDTO{}
	}
	request := RequestDTO{
		ParentCertificate: certificate.CertificateToDTO(&req.ParentCertificate),
		CertificateName:   req.CertificateName,
		CertificateType:   certificate.TypeToString(req.CertificateType),
		Subject:           user.ModelToUserBaseDTO(&req.Subject),
		Datetime:          req.Datetime,
		Status:            RequestStatusToString(req.Status),
	}
	return &request
}

func RequestStatusToString(reqStatus RequestStatus) string {
	switch reqStatus {
	case PENDING:
		return "PENDING"
	case ACCEPTED:
		return "ACCEPTED"
	case REJECTED:
		return "REJECTED"
	}
	// TODO error
	return ""
}

func StringToRequestStatus(certStatus string) RequestStatus {
	switch certStatus {
	case "PENDING":
		return PENDING
	case "ACCEPTED":
		return ACCEPTED
	case "REJECTED":
		return REJECTED
	}
	// TODO error
	return -1
}
