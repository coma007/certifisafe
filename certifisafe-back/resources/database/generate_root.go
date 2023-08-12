package database

import (
	"bytes"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func GenerateRoot(db certificate.CertificateRepository) error {
	config := utils.Config()
	// CA, root
	root := &x509.Certificate{
		Version:      3,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			CommonName:    config["name"],
			Organization:  []string{config["organization"]},
			Country:       []string{config["country"]},
			StreetAddress: []string{config["street"]},
			PostalCode:    []string{config["postal"]},
		},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// generate private key for CA (private key contains public)
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// create CA root certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, root, root, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	// create encoder
	rootPEM := new(bytes.Buffer)
	pem.Encode(rootPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// encode private key
	rootPrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(rootPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	serial := new(int64)
	*serial = root.SerialNumber.Int64()
	rootModel := &certificate.Certificate{
		//Id:        serial,
		Name:      root.Subject.CommonName,
		Issuer:    user.User{},
		Subject:   user.User{},
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Status:    certificate.CertificateStatus(certificate.ACTIVE),
		Type:      certificate.CertificateType(certificate.ROOT),
	}

	_, err = db.CreateCertificate(*rootModel)
	if err != nil {
		panic(err)
	}
	return nil
}
