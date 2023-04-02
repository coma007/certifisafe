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
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

func main() {
	yamlString, err := os.ReadFile("config.yaml")
	utils.CheckError(err)

	var config map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlString), &config)
	utils.CheckError(err)
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
	certificateService := service.NewDefaultCertificateService(certificateInMemoryRepository)
	certificateController := controller.NewCertificateHandler(certificateService)

	requestRepository := repository.NewRequestRepository(db, certificateInMemoryRepository)
	requestService := service.NewRequestServiceImpl(requestRepository, certificateService)
	requestController := controller.NewRequestController(requestService)

	router := httprouter.New()

	router.PATCH("/api/certificate/:id", certificateController.UpdateCertificate)
	router.GET("/api/certificate/:id", certificateController.GetCertificate)
	router.DELETE("/api/certificate/:id", certificateController.DeleteCertificate)
	router.POST("/api/certificate", certificateController.CreateCertificate)

	router.GET("/api/request", requestController.GetAllRequests)
	router.GET("/api/request/:id", requestController.GetRequest)
	router.POST("/api/request", requestController.CreateRequest)
	router.PATCH("/api/request/accept/:id", requestController.AcceptRequest)
	router.PATCH("/api/request/decline/:id", requestController.DeclineRequest)
	router.PATCH("/api/request/delete/:id", requestController.DeleteRequest)

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
