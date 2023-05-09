package main

import (
	"bufio"
	"bytes"
	"certifisafe-back/features/auth"
	certificate2 "certifisafe-back/features/certificate"
	"certifisafe-back/features/password_recovery"
	request2 "certifisafe-back/features/request"
	user2 "certifisafe-back/features/user"
	"certifisafe-back/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

var authService auth.AuthService

func main() {
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

	userRepository := user2.NewDefaultUserRepository(db)
	//passwordRecoveryRepository := password_recovery.NewDefaultPasswordRecoveryRepository(db)
	//verificationRepository := auth.NewInMemoryVerificationRepository(db)
	//
	//authService = auth.NewDefaultAuthService(userRepository, passwordRecoveryRepository, verificationRepository)
	//authController := auth.NewAuthHandler(authService)

	certificateRepository := certificate2.NewDefaultCertificateRepository(db)
	certificateFileStoreRepository := certificate2.NewDefaultFileStoreCertificateRepository()
	certificateService := certificate2.NewDefaultCertificateService(certificateRepository, certificateFileStoreRepository, userRepository)
	certificateController := certificate2.NewCertificateController(certificateService, authService)
	//
	//requestRepository := request2.NewDefaultRequestRepository(db, certificateRepository)
	//requestService := request2.NewDefaultRequestService(requestRepository, certificateService, userRepository)
	//requestController := request2.NewRequestController(requestService, certificateService, authService)

	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.HandleFunc("/api/certificate/:id", certificateController.GetCertificate).Methods("GET")
	router.HandleFunc("/api/certificate", certificateController.GetCertificates).Methods("GET")
	//router.GET("/api/certificate/:id/download", certificateController.DownloadCertificate)
	//router.PATCH("/api/certificate/:id/withdraw", certificateController.WithdrawCertificate)
	//router.GET("/api/certificate/:id/valid", certificateController.IsValid)
	//router.POST("/api/certificate/valid", certificateController.IsValidFile)
	//
	//router.POST("/api/request", requestController.CreateRequest)
	//router.GET("/api/request/:id", requestController.GetRequest)
	//router.GET("/api/request", requestController.GetAllRequestsByUser)
	//router.PATCH("/api/request/accept/:id", requestController.AcceptRequest)
	//router.PATCH("/api/request/decline/:id", requestController.DeclineRequest)
	//router.PATCH("/api/request/delete/:id", requestController.DeleteRequest)
	//router.POST("/api/certificate/generate", requestController.GenerateCertificates)
	//
	//router.POST("/api/login", authController.Login)
	//router.POST("/api/register", authController.Register)
	//router.GET("/api/verify-email/:verificationCode", authController.VerifyEmail)
	//router.POST("/api/password-recovery-request", authController.PasswordRecoveryRequest)
	//router.POST("/api/password-recovery", authController.PasswordRecovery)

	//router.HandlerFunc("GET", "/*any", corsMiddleware)
	//router.HandlerFunc("PATCH", "/*any", corsMiddleware)
	//router.HandlerFunc("POST", "/*any", corsMiddleware)
	//router.HandlerFunc("PUT", "/*any", corsMiddleware)
	//router.HandlerFunc("DELETE", "/*any", corsMiddleware)
	//createRoot(*certificateFileStoreRepository, certificateRepository)
	//runScript(db, "resources/database/data.sql")

	fmt.Println("http server runs on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func automigrate(db *gorm.DB) {
	err := db.AutoMigrate(&user2.User{}, &certificate2.Certificate{})
	utils.CheckError(err)
	err = db.AutoMigrate(&request2.Request{})
	utils.CheckError(err)
	err = db.AutoMigrate(&password_recovery.PasswordRecoveryRequest{})
	utils.CheckError(err)
	err = db.AutoMigrate(&auth.Verification{})
	utils.CheckError(err)
}

func createRoot(keyStore certificate2.DefaultFileStoreCertificateRepository, db certificate2.CertificateRepository) error {
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
	rootModel := &certificate2.Certificate{
		//Id:        serial,
		Name:      root.Subject.CommonName,
		Issuer:    user2.User{},
		Subject:   user2.User{},
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Status:    certificate2.CertificateStatus(certificate2.ACTIVE),
		Type:      certificate2.CertificateType(certificate2.ROOT),
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

func corsMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func middleware(n httprouter.Handle) httprouter.Handle {
	// TODO what is usage of this ?
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenValid, err := authService.ValidateToken(token)
		if err != nil || !tokenValid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		n(w, r, ps)
	}
}
