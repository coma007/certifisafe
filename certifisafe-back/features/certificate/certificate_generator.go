package certificate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func GenerateRootCa(subject pkix.Name, serial uint64) (x509.Certificate, bytes.Buffer, bytes.Buffer, error) {
	serialNumber := new(big.Int).SetUint64(serial)
	// CA, root
	ca := &x509.Certificate{
		Version:               3,
		SerialNumber:          serialNumber,
		Subject:               subject,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	// generate private key for CA (private key contains public)
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	// create CA root certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	// create encoder
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// encode private key
	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	return *ca, *caPEM, *caPrivKeyPEM, nil
}

func GenerateSubordinateCa(subject pkix.Name, serial uint64, rootTemplate *x509.Certificate, caPrivKey *rsa.PrivateKey) (x509.Certificate, bytes.Buffer, bytes.Buffer, error) {
	serialNumber := new(big.Int).SetUint64(serial)
	subject.SerialNumber = serialNumber.String()
	subTemplate := &x509.Certificate{
		Version:               3,
		SerialNumber:          serialNumber,
		Subject:               subject,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning, x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	//generate private key
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	// create certificate and sign it with CA key
	certBytes, err := x509.CreateCertificate(rand.Reader, subTemplate, rootTemplate, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	//create buffer and fill it with encoded value
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	//sign certificate with CA key
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return *subTemplate, *certPEM, *certPrivKeyPEM, nil
}

func GenerateLeafCert(subject pkix.Name, serial uint64, parent *x509.Certificate, parentPrivKey *rsa.PrivateKey) (x509.Certificate, bytes.Buffer, bytes.Buffer, error) {
	serialNumber := new(big.Int).SetUint64(serial)
	subject.SerialNumber = serialNumber.String()
	certTemplate := &x509.Certificate{
		Version:            3,
		SerialNumber:       serialNumber,
		Subject:            subject,
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		IPAddresses:        []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(10, 0, 0),
		SubjectKeyId:       []byte{1, 2, 3, 4, 6},

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		IsCA:        false,
	}

	//generate private key
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	// create certificate and sign it with CA key
	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, parent, &certPrivKey.PublicKey, parentPrivKey)
	if err != nil {
		return x509.Certificate{}, bytes.Buffer{}, bytes.Buffer{}, err
	}

	//create buffer and fill it with encoded value
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	//sign certificate with CA key
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return *certTemplate, *certPEM, *certPrivKeyPEM, nil
}
