package main

import (
	"certifisafe-back/controller"
	"certifisafe-back/domain"
	"certifisafe-back/repository"
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
	certificateService := domain.NewDefaultCertificateService(certificateInMemoryRepository)
	certificateController := controller.NewCertificateHandler(certificateService)
	fmt.Println(certificateController)

	router := httprouter.New()

	router.PATCH("/movies/:id", certificateController.UpdateCertificate)

	fmt.Println("http server runs on :8080")
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}
