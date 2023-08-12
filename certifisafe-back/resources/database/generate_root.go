package database

import (
	"certifisafe-back/features/certificate"
	"certifisafe-back/utils"
	"crypto/x509/pkix"
	"io/ioutil"
)

func GenerateRoot(db certificate.CertificateRepository) error {
	config := utils.Config()

	subject := pkix.Name{
		CommonName:    config["name"],
		Organization:  []string{config["organization"]},
		Country:       []string{config["country"]},
		StreetAddress: []string{config["street"]},
		PostalCode:    []string{config["postal"]},
	}

	_, caPEM, caPrivKeyPEM, err := certificate.GenerateRootCa(subject, 0)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("public/server.crt", caPEM.Bytes(), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("private/server.key", caPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		return err
	}

	return nil
}
