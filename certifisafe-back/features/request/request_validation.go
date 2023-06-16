package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

func (req NewRequestDTO) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CertificateName, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.ParentSerial, validation.Required),
		validation.Field(&req.CertificateType, validation.Required),
	)
}
