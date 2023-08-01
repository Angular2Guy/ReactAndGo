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

	database.DB.Migrator().AddColumn(&aumodel.PostCodeLocation{}, "State")
	database.DB.Migrator().AddColumn(&aumodel.PostCodeLocation{}, "County")

	database.DB.Migrator().AddColumn(&gsmodel.GasStation{}, "State")
	database.DB.Migrator().AddColumn(&gsmodel.GasStation{}, "County")

	log.Printf("DB Migration Done.")
}
