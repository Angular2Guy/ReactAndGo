package appuser

import (
	"angular-and-go/pkd/appuser/aumodel"
	"angular-and-go/pkd/database"
	"log"

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

const (
	Ok SigninResult = iota
	UsernameTaken
	Invalid
	Failed
)

func Login(appUserIn AppUserIn) string {
	result := ""
	var appUser aumodel.AppUser
	if err := database.DB.Where("username = ?", appUserIn.Username).First(&appUser); err != nil {
		log.Printf("User not found: %v", appUserIn.Username)
		return result
	}
	if err := bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(appUserIn.Password)); err != nil {
		log.Printf("Password wrong user: %v", appUser.Username)
		return result
	}
	// add jwt token creation
	return result
}

func Signin(appUserIn AppUserIn) SigninResult {
	var result SigninResult = Invalid
	if len(appUserIn.Username) < 4 || len(appUserIn.Password) < 8 {
		return result
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var appUser aumodel.AppUser
		if err := tx.Where("name = ?", appUserIn.Username).First(&appUser); err == nil {
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

func generatePasswordHash(password string) []byte {
	passwordSlice := []byte(password)
	hashValue, err := bcrypt.GenerateFromPassword(passwordSlice, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("GenerateFromPassword failed: %v\n", err.Error())
	}
	return hashValue
}