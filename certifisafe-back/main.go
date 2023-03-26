package main

import (
	"certifisafe-back/controller"
	"certifisafe-back/repository"
	"certifisafe-back/service"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CertifiSafe")
}

func main() {
	certificateInMemoryRepository := repository.NewInMemoryCertificateRepository()
	certificateService := service.NewDefaultCertificateService(certificateInMemoryRepository)
	certificateController := controller.NewCertificateHandler(certificateService)
	fmt.Println(certificateController)

	router := httprouter.New()

	router.PATCH("/movies/:id", certificateController.UpdateCertificate)

	fmt.Println("http server runs on :8080")
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}

//func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Verify the user's credentials
//		if !verifyCredentials(r) {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		// If the credentials are valid, allow the request to proceed
//		next(w, r)
//	}
//}
