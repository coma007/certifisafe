package main

import (
	"bytes"
	"certifisafe-back/controller"
	"certifisafe-back/model"
	"certifisafe-back/repository"
	"certifisafe-back/service"
	"certifisafe-back/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

var auth service.IAuthService

func main() {
	config := utils.Config()
	password := config["password"]
	user := config["user"]

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:5432/certifisafe?sslmode=disable", user, password))
	utils.CheckError(err)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	userInMemoryRepository := repository.NewInMemoryUserRepository(db)
	passwordRecoveryInMemoryRepository := repository.NewInMemoryPasswordRecoveryRepository(db)
	verificationInMemoryRepository := repository.NewInMemoryVerificationRepository(db)
	auth = service.NewAuthService(userInMemoryRepository, passwordRecoveryInMemoryRepository, verificationInMemoryRepository)
	authController := controller.NewAuthHandler(auth)

	certificateInMemoryRepository := repository.NewInMemoryCertificateRepository(db)
	certificateKeyStoreInMemoryRepository := repository.NewInMemoryCertificateKeyStoreRepository(db)
	certificateService := service.NewDefaultCertificateService(certificateInMemoryRepository, certificateKeyStoreInMemoryRepository)
	certificateController := controller.NewCertificateHandler(certificateService)

	requestRepository := repository.NewRequestRepository(db, certificateInMemoryRepository)
	requestService := service.NewRequestServiceImpl(requestRepository, certificateService)
	requestController := controller.NewRequestController(requestService, auth)

	router := httprouter.New()

	router.GET("/api/certificate/:id", certificateController.GetCertificate)
	router.GET("/api/certificate", certificateController.GetCertificates)
	router.DELETE("/api/certificate/:id", certificateController.DeleteCertificate)
	router.POST("/api/certificate", certificateController.CreateCertificate)
	router.POST("/api/certificate/generate", certificateController.Generate)
	router.GET("/api/certificate/:id/valid", certificateController.IsValid)

	router.GET("/api/request", requestController.GetAllRequests)
	router.GET("/api/request/:id", requestController.GetRequest)
	router.POST("/api/request", requestController.CreateRequest)
	router.PATCH("/api/request/accept/:id", requestController.AcceptRequest)
	router.PATCH("/api/request/decline/:id", requestController.DeclineRequest)
	router.PATCH("/api/request/delete/:id", requestController.DeleteRequest)

	router.POST("/api/login", authController.Login)
	router.POST("/api/register", authController.Register)
	router.POST("/api/password-recovery-request", authController.PasswordRecoveryRequest)
	router.POST("/api/password-recovery", authController.PasswordRecovery)

	runScript(db, "utils/schema.sql")
	createRoot(*certificateKeyStoreInMemoryRepository, certificateInMemoryRepository)

	runScript(db, "utils/data.sql")

	fmt.Println("http server runs on :8080")
	err = http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

func createRoot(keyStore repository.InmemoryKeyStoreCertificateRepository, db repository.ICertificateRepository) error {
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

	rootModel := &model.Certificate{
		Id:        root.SerialNumber.Int64(),
		Name:      root.Subject.CommonName,
		Issuer:    nil,
		Subject:   nil,
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Status:    model.CertificateStatus(model.ACTIVE),
		Type:      model.CertificateType(model.ROOT),
		PublicKey: 0,
	}

	_, err = keyStore.CreateCertificate(*root.SerialNumber, *rootPEM, *rootPrivateKeyPEM)
	_, err = db.CreateCertificate(*rootModel)
	if err != nil {
		panic(err)
	}
	return nil
}

func runScript(db *sql.DB, script string) {
	c, ioErr := os.ReadFile(script)
	utils.CheckError(ioErr)
	commands := string(c)
	_, err := db.Exec(commands)
	if err != nil {
		//panic("Couldn't load sql script")
		panic(err)
	}
}

func middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenValid, err := auth.ValidateToken(token)
		if err != nil || !tokenValid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		n(w, r, ps)
	}
}
