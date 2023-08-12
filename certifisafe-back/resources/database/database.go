package database

import (
	"bufio"
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/request"
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewDatabase() *gorm.DB {
	config := utils.Config()
	password := config["password"]
	dbuser := config["user"]

	dbString := fmt.Sprintf("postgres://%s:%s@localhost:5432/certifisafe?sslmode=disabled",
		dbuser, password, "root.crt", "srv.crt", "srv.key")
	dbPostgree := postgres.Open(dbString)

	db, err := gorm.Open(dbPostgree, &gorm.Config{PrepareStmt: true, TranslateError: true})
	automigrate(db)
	utils.CheckError(err)
	return db
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
