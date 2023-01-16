package database

import (
	"angular-and-go/pkd/appuser/aumodel"
	"angular-and-go/pkd/gasstation/gsmodel"
	"log"
)

func MigrateDB() {
	//DB.AutoMigrate(&gsmodel.GasStation{})
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
	log.Printf("DB Migration Done.")
}
