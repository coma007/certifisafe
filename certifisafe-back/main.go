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

	userInMemoryRepository := repository.NewInMemoryUserRepository(db)
	auth = service.NewAuthService(userInMemoryRepository)
	authController := controller.NewAuthHandler(auth)

	fmt.Println(certificateController)

	router := httprouter.New()

	router.PATCH("/certificate/:id", certificateController.UpdateCertificate)
	router.GET("/certificate/:id", middleware(certificateController.GetCertificate))
	router.DELETE("/certificate/:id", certificateController.DeleteCertificate)
	router.POST("/certificate", certificateController.CreateCertificate)

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
		panic("Couldn't load sql script")
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
