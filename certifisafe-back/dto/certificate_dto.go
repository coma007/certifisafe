package dto

import (
	"certifisafe-back/model"
	"crypto/x509"
	"time"
)

type CertificateDTO struct {
	Serial      *uint64
	Name        string
	ValidFrom   time.Time
	ValidTo     time.Time
	IssuerName  string
	SubjectName string
	Status      string
	Type        string
}

func CertificateToDTO(cert *model.Certificate) *CertificateDTO {
	if cert == nil {
		return nil
	}
	serial := uint64(cert.ID)
	certificate := CertificateDTO{
		Serial:    &serial,
		Name:      cert.Name,
		ValidFrom: cert.ValidFrom,
		ValidTo:   cert.ValidTo,
		// TODO make nested object user
		IssuerName:  cert.Issuer.FirstName + " " + cert.Issuer.LastName,
		SubjectName: cert.Subject.FirstName + " " + cert.Subject.LastName,
		Status:      TypeToString(cert.Type),
		Type:        StatusToString(cert.Status),
	}
	return &certificate
}

func CertificateDTOtoModel(cert *CertificateDTO) *model.Certificate {
	if cert == nil {
		return nil
	}
	certificate := model.Certificate{
		//ID:        cert.Serial,
		Issuer: model.User{
			Email:     "",
			Password:  "",
			FirstName: "",
			LastName:  "",
			Phone:     "",
			IsAdmin:   false,
		},
		Subject:   model.User{},
		ValidFrom: cert.ValidFrom,
		ValidTo:   cert.ValidTo,
		Status:    StringToStatus(cert.Status),
		Type:      StringToType(cert.Type),
	}
	return &certificate
}

func X509CertificateToCertificateDTO(cert *x509.Certificate) *CertificateDTO {
	if cert == nil {
		return nil
	}
	serial := cert.SerialNumber.Uint64()
	certificate := CertificateDTO{
		Serial:    &serial,
		Name:      cert.Subject.CommonName,
		ValidFrom: cert.NotBefore,
		ValidTo:   cert.NotAfter,
		// TODO make nested object user
		IssuerName:  cert.Issuer.CommonName,
		SubjectName: cert.Subject.CommonName,
	}
	return &certificate
}

func TypeToString(certType model.CertificateType) string {
	switch certType {
	case model.ROOT:
		return "ROOT"
	case model.INTERMEDIATE:
		return "INTERMEDIATE"
	case model.END:
		return "END"
	}
	// TODO error
	return ""
}

func StringToType(certType string) model.CertificateType {
	switch certType {
	case "ROOT":
		return model.ROOT
	case "INTERMEDIATE":
		return model.INTERMEDIATE
	case "END":
		return model.END
	}
	// TODO error
	return -1
}

func StatusToString(certStatus model.CertificateStatus) string {
	switch certStatus {
	case model.ACTIVE:
		return "ACTIVE"
	case model.EXPIRED:
		return "EXPIRED"
	case model.WITHDRAWN:
		return "WITHDRAWN"
	}
	// TODO error
	return ""
}

func StringToStatus(certStatus string) model.CertificateStatus {
	switch certStatus {
	case "ACTIVE":
		return model.ACTIVE
	case "EXPIRED":
		return model.EXPIRED
	case "WITHDRAWN":
		return model.WITHDRAWN
	}
	// TODO error
	return -1
}
