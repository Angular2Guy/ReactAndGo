package database

import (
	"log"
	"react-and-go/pkd/appuser/aumodel"
	"react-and-go/pkd/gasstation/gsmodel"
)

func MigrateDB() {
	if !DB.Migrator().HasTable(&gsmodel.GasStation{}) {
		DB.AutoMigrate(&gsmodel.GasStation{})
	}
	if !DB.Migrator().HasTable(&gsmodel.GasPrice{}) {
		DB.AutoMigrate(&gsmodel.GasPrice{})
	}
	if !DB.Migrator().HasTable(&aumodel.AppUser{}) {
		DB.AutoMigrate(&aumodel.AppUser{})
	}
	if !DB.Migrator().HasTable(&aumodel.LoggedOutUser{}) {
		DB.AutoMigrate(&aumodel.LoggedOutUser{})
	}
	if !DB.Migrator().HasTable(&aumodel.PostCodeLocation{}) {
		DB.AutoMigrate(&aumodel.PostCodeLocation{})
	}
	log.Printf("DB Migration Done.")
}
