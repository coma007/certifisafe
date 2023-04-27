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
	"encoding/pem"
	"fmt"
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

var auth service.IAuthService

func main() {
	config := utils.Config()
	password := config["password"]
	user := config["user"]

	dbPostgree := postgres.Open(fmt.Sprintf("postgres://%s:%s@localhost:5432/certifisafe?sslmode=disable", user, password))
	db, err := gorm.Open(dbPostgree, &gorm.Config{PrepareStmt: true})
	//DeleteCreatedEntities(db)
	//runScript(db, "utils/schema.sql")
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
	//utils.CheckError(err)
	//
	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(db)

	userInMemoryRepository := repository.NewInMemoryUserRepository(db)
	passwordRecoveryInMemoryRepository := repository.NewInMemoryPasswordRecoveryRepository(db)
	auth = service.NewAuthService(userInMemoryRepository, passwordRecoveryInMemoryRepository)
	authController := controller.NewAuthHandler(auth)

	certificateInMemoryRepository := repository.NewInMemoryCertificateRepository(db)
	certificateKeyStoreInMemoryRepository := repository.NewInMemoryCertificateKeyStoreRepository()
	certificateService := service.NewDefaultCertificateService(certificateInMemoryRepository, certificateKeyStoreInMemoryRepository, userInMemoryRepository)
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

	//
	//createRoot(*certificateKeyStoreInMemoryRepository, certificateInMemoryRepository)

	//runScript(db, "utils/data.sql")

	fmt.Println("http server runs on :8080")
	err = http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

func automigrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Certificate{})
	utils.CheckError(err)
	err = db.AutoMigrate(&model.Request{})
	utils.CheckError(err)
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

	serial := new(int64)
	*serial = root.SerialNumber.Int64()
	rootModel := &model.Certificate{
		Id:        serial,
		Name:      root.Subject.CommonName,
		Issuer:    model.User{},
		Subject:   model.User{},
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Status:    model.CertificateStatus(model.ACTIVE),
		Type:      model.CertificateType(model.ROOT),
	}

	_, err = keyStore.CreateCertificate(root.SerialNumber.Int64(), *rootPEM, *rootPrivateKeyPEM)
	_, err = db.CreateCertificate(*rootModel)
	if err != nil {
		panic(err)
	}
	return nil
}

func runScript(db *gorm.DB, script string) {
	c, ioErr := os.ReadFile(script)
	utils.CheckError(ioErr)
	commands := string(c)
	result := db.Raw(commands)
	if result.Error != nil {
		//panic("Couldn't load sql script")
		panic(result.Error)
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

func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	// Remove the hook once we're done
	defer db.Callback().Create().Remove(hookName)
	// Find out if the current db object is already a transaction
	tx := db
	tx = db.Begin()
	// Loop from the end. It is important that we delete the entries in the
	// reverse order of their insertion
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
		tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
	}
	tx.Commit()
	return nil
}
