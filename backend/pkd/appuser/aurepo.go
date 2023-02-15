package appuser

import (
	"angular-and-go/pkd/appuser/aumodel"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/token"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AppUserIn struct {
	Username  string
	Password  string
	Uuid      string
	Latitude  float64
	Longitude float64
}

type SigninResult int

type PostCodeData struct {
	Label           string
	PostCode        int32
	Population      int32
	SquareKM        float32
	CenterLongitude float64
	CenterLatitude  float64
}

const (
	Ok SigninResult = iota
	UsernameTaken
	Invalid
	Failed
)

func FindLocation(locationStr string) []aumodel.PostCodeLocation {
	result := []aumodel.PostCodeLocation{}
	database.DB.Where("lower(label) like ?", fmt.Sprintf("%%%v%%", strings.ToLower(strings.TrimSpace(locationStr)))).Limit(20).Find(&result)
	//log.Printf("Select: %v failed. %v", fmt.Sprintf("%%%v%%", strings.ToLower(strings.TrimSpace(locationStr))), err)
	return result
}

func Login(appUserIn AppUserIn) (string, int) {
	result := ""
	status := http.StatusUnauthorized
	//log.Printf("%v", appUserIn.Username)
	var appUser aumodel.AppUser
	if err := database.DB.Where("username = ?", appUserIn.Username).First(&appUser); err.Error != nil {
		log.Printf("User not found: %v error: %v\n", appUserIn.Username, err.Error)
		return result, status
	}
	if err := bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(appUserIn.Password)); err != nil {
		log.Printf("Password wrong. Username: %v\n", appUser.Username)
		return result, status
	}
	//jwt token creation
	result, err := token.CreateToken(token.TokenUser{Username: appUser.Username, Roles: []string{"USERS"}})
	if err != nil {
		log.Printf("Failed to create jwt token: %v\n", err)
	} else {
		status = http.StatusOK
	}
	return result, status
}

func Signin(appUserIn AppUserIn) SigninResult {
	var result SigninResult = Invalid
	if len(appUserIn.Username) < 4 || len(appUserIn.Password) < 8 {
		return result
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var appUser aumodel.AppUser
		if err := tx.Where("username = ?", appUserIn.Username).First(&appUser); err.Error == nil {
			result = UsernameTaken
			return nil
		}
		appUser.Username = appUserIn.Username
		appUser.Password = string(generatePasswordHash(appUserIn.Password))
		appUser.Latitude = appUserIn.Latitude
		appUser.Longitude = appUserIn.Longitude
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

func ImportPostCodeData(postCodeData []PostCodeData) {
	postCodeLocations := mapToPostCodeLocation(postCodeData)
	var oriPostCodeLocations []aumodel.PostCodeLocation
	database.DB.Find(&oriPostCodeLocations)
	postCodeLocationsMap := make(map[int32]aumodel.PostCodeLocation)
	for _, oriPostCodeLocation := range oriPostCodeLocations {
		postCodeLocationsMap[oriPostCodeLocation.PostCode] = oriPostCodeLocation
	}

	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, postCodeLocation := range postCodeLocations {
			oriPostCodeLocation, exists := postCodeLocationsMap[postCodeLocation.PostCode]
			if exists {
				oriPostCodeLocation.Label = postCodeLocation.Label
				oriPostCodeLocation.PostCode = postCodeLocation.PostCode
				oriPostCodeLocation.Population = postCodeLocation.Population
				oriPostCodeLocation.SquareKM = postCodeLocation.SquareKM
				oriPostCodeLocation.CenterLongitude = postCodeLocation.CenterLongitude
				oriPostCodeLocation.CenterLatitude = postCodeLocation.CenterLatitude
				tx.Save(&oriPostCodeLocation)
			} else {
				tx.Save(&postCodeLocation)
			}
		}
		return nil
	})
	log.Printf("PostCodeLocations saved: %v\n", len(postCodeLocations))
}

func mapToPostCodeLocation(postCodeData []PostCodeData) []aumodel.PostCodeLocation {
	result := []aumodel.PostCodeLocation{}
	for _, myPostCodeData := range postCodeData {
		myPostCodeLocation := aumodel.PostCodeLocation{}
		myPostCodeLocation.Label = myPostCodeData.Label
		myPostCodeLocation.PostCode = myPostCodeData.PostCode
		myPostCodeLocation.Population = myPostCodeData.Population
		myPostCodeLocation.SquareKM = myPostCodeData.SquareKM
		myPostCodeLocation.CenterLongitude = myPostCodeData.CenterLongitude
		myPostCodeLocation.CenterLatitude = myPostCodeData.CenterLatitude
		result = append(result, myPostCodeLocation)
	}
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
