package main

import (
	"certifisafe-back/controller"
	"certifisafe-back/repository"
	"certifisafe-back/service"
	"certifisafe-back/utils"
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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

	runScript(db)

	certificateInMemoryRepository := repository.NewInMemoryCertificateRepository(db)
	certificateKeyStoreInMemoryRepository := repository.NewInMemoryCertificateKeyStoreRepository(db)
	certificateService := service.NewDefaultCertificateService(certificateInMemoryRepository, certificateKeyStoreInMemoryRepository)
	certificateController := controller.NewCertificateHandler(certificateService)

	requestRepository := repository.NewRequestRepository(db, certificateInMemoryRepository)
	requestService := service.NewRequestServiceImpl(requestRepository, certificateService)
	requestController := controller.NewRequestController(requestService)

	userInMemoryRepository := repository.NewInMemoryUserRepository(db)
	auth = service.NewAuthService(userInMemoryRepository)
	authController := controller.NewAuthHandler(auth)

	router := httprouter.New()

	router.GET("/api/certificate/:id", certificateController.GetCertificate)
	router.DELETE("/api/certificate/:id", certificateController.DeleteCertificate)
	router.POST("/api/certificate", certificateController.CreateCertificate)
	router.GET("/api/certificate/:id/valid", certificateController.IsValid)

	router.GET("/api/request", requestController.GetAllRequests)
	router.GET("/api/request/:id", requestController.GetRequest)
	router.POST("/api/request", requestController.CreateRequest)
	router.PATCH("/api/request/accept/:id", requestController.AcceptRequest)
	router.PATCH("/api/request/decline/:id", requestController.DeclineRequest)
	router.PATCH("/api/request/delete/:id", requestController.DeleteRequest)

	router.POST("/login", authController.Login)
	router.GET("/validate", authController.Validate)

	fmt.Println("http server runs on :8080")
	err = http.ListenAndServe(":8080", router)
	log.Fatal(err)

}

func runScript(db *sql.DB) {
	c, ioErr := os.ReadFile("utils/data.sql")
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
