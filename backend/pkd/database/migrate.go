package database

import "angular-and-go/pkd/gasstation/gsmodel"

//import "angular-and-go/pkd/gasstation/gsmodel"

func MigrateDB() {
	//DB.AutoMigrate(&gsmodel.GasStation{})
	if !DB.Migrator().HasTable(&gsmodel.GasStation{}) {
		DB.AutoMigrate(&gsmodel.GasStation{})
	}
	if !DB.Migrator().HasTable(&gsmodel.GasPrice{}) {
		DB.AutoMigrate(&gsmodel.GasPrice{})
	}

}
