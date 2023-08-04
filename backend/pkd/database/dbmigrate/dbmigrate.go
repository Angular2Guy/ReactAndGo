/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package dbmigrate

import (
	"log"
	"react-and-go/pkd/appuser/aumodel"
	pcmodel "react-and-go/pkd/appuser/pcmodel"
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
	if !database.DB.Migrator().HasTable(&pcmodel.PostCodeLocation{}) {
		database.DB.AutoMigrate(&pcmodel.PostCodeLocation{})
	}
	if !database.DB.Migrator().HasTable(&unmodel.UserNotification{}) {
		database.DB.AutoMigrate(&unmodel.UserNotification{})
	}
	if !database.DB.Migrator().HasColumn(&pcmodel.PostCodeLocation{}, "StateDataID") {
		database.DB.Migrator().AddColumn(&pcmodel.PostCodeLocation{}, "StateDataID")
	}
	if !database.DB.Migrator().HasColumn(&pcmodel.PostCodeLocation{}, "CountyDataID") {
		database.DB.Migrator().AddColumn(&pcmodel.PostCodeLocation{}, "CountyDataID")
	}
	if !database.DB.Migrator().HasTable(&pcmodel.CountyData{}) {
		database.DB.AutoMigrate(&pcmodel.CountyData{})
	}
	if !database.DB.Migrator().HasTable(&pcmodel.StateData{}) {
		database.DB.AutoMigrate(&pcmodel.StateData{})
	}
	log.Printf("DB Migration Done.")
}
