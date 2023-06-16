package main

import (
	"bufio"
	"bytes"
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/request"
	"certifisafe-back/features/user"
	"certifisafe-back/internal"
	"certifisafe-back/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("\033[33mThis is yellow text\033[0m")
	fmt.Println("\033[32mThis is green text\033[0m")
	fmt.Println("\033[31mThis is red text\033[0m")
	config := utils.Config()
	password := config["password"]
	dbuser := config["user"]

	dbPostgree := postgres.Open(fmt.Sprintf("postgres://%s:%s@localhost:5432/certifisafe?sslmode=disable", dbuser, password))
	db, err := gorm.Open(dbPostgree, &gorm.Config{PrepareStmt: true, TranslateError: true})
	automigrate(db)
	utils.CheckError(err)

	defer func(db *gorm.DB) {
		sqlDb, err := db.DB()
		if err != nil {
			panic(err)
		}
		err = sqlDb.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	app := internal.NewDefaultAppFactory(db)
	app.InitApp()

	router := internal.NewDefaultRouter(app)
	router.ListenAndServe()
}

func automigrate(db *gorm.DB) {
	err := db.AutoMigrate(&user.User{}, &certificate.Certificate{})
	utils.CheckError(err)
	err = db.AutoMigrate(&request.Request{})
	utils.CheckError(err)
	err = db.AutoMigrate(&password_recovery.PasswordRecoveryRequest{})
	utils.CheckError(err)
	err = db.AutoMigrate(&auth.Verification{})
	utils.CheckError(err)
	err = db.AutoMigrate(&password_recovery.PasswordHistory{})
	utils.CheckError(err)
}

func createRoot(keyStore certificate.DefaultFileStoreCertificateRepository, db certificate.CertificateRepository) error {
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

	//_, err = keyStore.CreateCertificate(root.SerialNumber.Int64(), *rootPEM, *rootPrivateKeyPEM)
	_, err = db.CreateCertificate(*rootModel)
	if err != nil {
		panic(err)
	}
	return nil
}

func runScript(db *gorm.DB, script string) {
	file, err := os.Open(script)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		db.Exec(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
