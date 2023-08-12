package main

import (
	"certifisafe-back/internal"
	"certifisafe-back/resources/database"
	"gorm.io/gorm"
)

func main() {

	db := database.NewDatabase()

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

	app := internal.NewDefaultAppFactory(db)
	app.InitApp()

	// if needed, uncomment this:
	database.GenerateRoot()

	router := internal.NewDefaultRouter(app)
	router.ListenAndServe()

}
