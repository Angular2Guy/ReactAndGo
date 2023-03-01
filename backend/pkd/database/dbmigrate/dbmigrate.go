package dbmigrate

import (
	"log"
	"react-and-go/pkd/appuser/aumodel"
	database "react-and-go/pkd/database"
	"react-and-go/pkd/gasstation/gsmodel"
	unmodel "react-and-go/pkd/notification/model"
)

func MigrateDB() {
	if !database.DB.Migrator().HasTable(&gsmodel.GasStation{}) {
		database.DB.AutoMigrate(&gsmodel.GasStation{})
	}
	if !database.DB.Migrator().HasTable(&gsmodel.GasPrice{}) {
		database.DB.AutoMigrate(&gsmodel.GasPrice{})
	}
	if !database.DB.Migrator().HasTable(&aumodel.AppUser{}) {
		database.DB.AutoMigrate(&aumodel.AppUser{})
	}
	if !database.DB.Migrator().HasTable(&aumodel.LoggedOutUser{}) {
		database.DB.AutoMigrate(&aumodel.LoggedOutUser{})
	}
	if !database.DB.Migrator().HasTable(&aumodel.PostCodeLocation{}) {
		database.DB.AutoMigrate(&aumodel.PostCodeLocation{})
	}
	if !database.DB.Migrator().HasTable(&unmodel.UserNotification{}) {
		database.DB.AutoMigrate(&unmodel.UserNotification{})
	}
	log.Printf("DB Migration Done.")
}
