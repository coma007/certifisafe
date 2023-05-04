package certificate

import (
	"certifisafe-back/features/user"
	"time"
)

// TODO add organization, postal code, ... to Certificate

type CertificateDTO struct {
	Serial    *uint64
	Name      string
	ValidFrom time.Time
	ValidTo   time.Time
	Issuer    user.UserBaseDTO
	Subject   user.UserBaseDTO
	Status    string
	Type      string
}

func CertificateToDTO(cert *Certificate) *CertificateDTO {
	if cert == nil {
		return nil
	}
	serial := uint64(cert.ID)
	certificate := CertificateDTO{
		Serial:    &serial,
		Name:      cert.Name,
		ValidFrom: cert.ValidFrom,
		ValidTo:   cert.ValidTo,
		Issuer:    *user.ModelToUserBaseDTO(&cert.Issuer),
		Subject:   *user.ModelToUserBaseDTO(&cert.Issuer),
		Status:    TypeToString(cert.Type),
		Type:      StatusToString(cert.Status),
	}
	return &certificate
}

func ModelToCertificateDTO(cert *Certificate) *CertificateDTO {
	if cert == nil {
		return nil
	}
	serial := uint64(cert.ID)
	certificate := CertificateDTO{
		Serial:    &serial,
		Name:      cert.Name,
		ValidFrom: cert.ValidFrom,
		ValidTo:   cert.ValidTo,
		Issuer:    *user.ModelToUserBaseDTO(&cert.Issuer),
		Subject:   *user.ModelToUserBaseDTO(&cert.Subject),
	}
	return &certificate
}

func TypeToString(certType CertificateType) string {
	switch certType {
	case ROOT:
		return "ROOT"
	case INTERMEDIATE:
		return "INTERMEDIATE"
	case END:
		return "END"
	}
	// TODO error
	return ""
}

func StringToType(certType string) CertificateType {
	switch certType {
	case "ROOT":
		return ROOT
	case "INTERMEDIATE":
		return INTERMEDIATE
	case "END":
		return END
	}
	// TODO error
	return -1
}

func StatusToString(certStatus CertificateStatus) string {
	switch certStatus {
	case ACTIVE:
		return "ACTIVE"
	case EXPIRED:
		return "EXPIRED"
	case WITHDRAWN:
		return "WITHDRAWN"
	}
	// TODO error
	return ""
}

func StringToStatus(certStatus string) CertificateStatus {
	switch certStatus {
	case "ACTIVE":
		return ACTIVE
	case "EXPIRED":
		return EXPIRED
	case "WITHDRAWN":
		return WITHDRAWN
	}
	// TODO error
	return -1
}
