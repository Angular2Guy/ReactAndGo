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
package appuser

import (
	"log"
	"net/http"
	aumodel "react-and-go/pkd/appuser/aumodel"
	"react-and-go/pkd/database"
	token "react-and-go/pkd/token"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AppUserIn struct {
	Username     string
	Password     string
	Uuid         string
	Language     UserLang
	Latitude     float64
	Longitude    float64
	SearchRadius float64
}

type AppTargetIn struct {
	Username     string
	TargetDiesel string
	TargetE10    string
	TargetE5     string
}

type DbResult int

type UserLang string

const (
	German  UserLang = "de"
	English          = "en"
)

const (
	Ok DbResult = iota
	UsernameTaken
	Invalid
	Failed
)

func FindAllUsers() []aumodel.AppUser {
	var result []aumodel.AppUser
	database.DB.Find(&result)
	return result
}

func StoreUserLogout(username string, uuid string) []token.LoggedOutUserOut {
	database.DB.Delete("last_logout < ?", time.Now().Add(-4*time.Minute))
	loggedOutUser := aumodel.LoggedOutUser{Username: username, Uuid: uuid, LastLogout: time.Now()}
	database.DB.Save(loggedOutUser)
	var loggedOutUsers []aumodel.LoggedOutUser
	database.DB.Find(&loggedOutUsers)
	var results []token.LoggedOutUserOut
	for _, myLoggedOutUser := range loggedOutUsers {
		loggedOutUserOut := token.LoggedOutUserOut{Username: myLoggedOutUser.Username, Uuid: myLoggedOutUser.Uuid, LastLogout: myLoggedOutUser.LastLogout}
		results = append(results, loggedOutUserOut)
	}
	return results
}

func Login(appUserIn AppUserIn) (string, int, string, float64, float64, float64, int, int, int) {
	result := ""
	status := http.StatusUnauthorized
	//log.Printf("%v", appUserIn.Username)
	var appUser aumodel.AppUser
	if err := database.DB.Where("username = ?", appUserIn.Username).First(&appUser); err.Error != nil {
		log.Printf("User not found: %v error: %v\n", appUserIn.Username, err.Error)
		return result, status, "", 0.0, 0.0, 0.0, 0, 0, 0
	}
	if err := bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(appUserIn.Password)); err != nil {
		log.Printf("Password wrong. Username: %v\n", appUser.Username)
		return result, status, "", 0.0, 0.0, 0.0, 0, 0, 0
	}
	//jwt token creation
	result, err := token.CreateToken(token.TokenUser{Username: appUser.Username, Roles: []string{"USERS"}})
	if err != nil {
		log.Printf("Failed to create jwt token: %v\n", err)
		return result, status, "", 0.0, 0.0, 0.0, 0, 0, 0
	} else {
		status = http.StatusOK
	}
	return result, status, appUser.Uuid, appUser.Longitude, appUser.Latitude, appUser.SearchRadius, appUser.TargetE5, appUser.TargetE10, appUser.TargetDiesel
}

func Signin(appUserIn AppUserIn) DbResult {
	var result DbResult = Invalid
	if len(appUserIn.Username) < 4 || len(appUserIn.Password) < 8 {
		return result
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var appUser aumodel.AppUser
		//check usernames
		if err := tx.Where("username = ?", appUserIn.Username).First(&appUser); err.Error == nil {
			result = UsernameTaken
			return nil
		}
		//generate uuid
		myUuid, err := uuid.NewRandom()
		if err != nil {
			result = Failed
			return nil
		}
		appUser.Username = appUserIn.Username
		appUser.Password = string(generatePasswordHash(appUserIn.Password))
		appUser.Uuid = myUuid.String()
		appUser.LangKey = string(appUserIn.Language)
		tx.Save(&appUser)
		return nil
	})
	if err != nil {
		result = Failed
	} else if result == Invalid {
		result = Ok
	}
	return result
}

func StoreLocationAndRadius(appUserIn AppUserIn) DbResult {
	result := Invalid
	database.DB.Transaction(func(tx *gorm.DB) error {
		var appUser aumodel.AppUser
		if err := tx.Where("username = ?", appUserIn.Username).Find(&appUser); err.Error == nil {
			appUser.Longitude = appUserIn.Longitude
			appUser.Latitude = appUserIn.Latitude
			appUser.SearchRadius = appUserIn.SearchRadius
			tx.Save(&appUser)
			result = Ok
		}
		return nil
	})
	return result
}

func StoreTargetPrices(appTargetIn AppTargetIn) DbResult {
	result := Invalid
	database.DB.Transaction(func(tx *gorm.DB) error {
		var appUser aumodel.AppUser
		var txError error = nil
		if err := tx.Where("username = ?", appTargetIn.Username).Find(&appUser); err.Error == nil {
			if targetPrice, err := strconv.ParseInt(strings.ReplaceAll(appTargetIn.TargetDiesel, ".", ""), 10, 32); err == nil {
				appUser.TargetDiesel = int(targetPrice)
			} else {
				log.Printf("TargetDiesel: %v\n", appTargetIn.TargetDiesel)
				txError = err
			}
			if targetPrice, err := strconv.ParseInt(strings.ReplaceAll(appTargetIn.TargetE10, ".", ""), 10, 32); err == nil {
				appUser.TargetE10 = int(targetPrice)
			} else {
				log.Printf("TargetE10: %v\n", appTargetIn.TargetE10)
				txError = err
			}
			if targetPrice, err := strconv.ParseInt(strings.ReplaceAll(appTargetIn.TargetE5, ".", ""), 10, 32); err == nil {
				appUser.TargetE5 = int(targetPrice)
			} else {
				log.Printf("TargetE5: %v\n", appTargetIn.TargetE5)
				txError = err
			}
			if txError == nil {
				tx.Save(&appUser)
				result = Ok
			}
		}
		return txError
	})
	return result
}

func generatePasswordHash(password string) []byte {
	passwordSlice := []byte(password)
	hashValue, err := bcrypt.GenerateFromPassword(passwordSlice, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("GenerateFromPassword failed: %v\n", err.Error())
	}
	return hashValue
}
