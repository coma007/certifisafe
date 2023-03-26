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

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/certisafe?sslmode=disable", user, password))
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

	fmt.Println(certificateController)

	router := httprouter.New()

	router.PATCH("/certificate/:id", certificateController.UpdateCertificate)
	router.GET("/certificate/:id", certificateController.GetCertificate)
	router.DELETE("/certificate/:id", certificateController.DeleteCertificate)
	router.POST("/certificate", certificateController.CreateCertificate)

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
		panic("Couldn't load sql script")
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
