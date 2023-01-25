package appuser

import (
	"angular-and-go/pkd/appuser/aumodel"
	"angular-and-go/pkd/database"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

func Login(appUserIn AppUserIn) (string, int) {
	result := ""
	status := http.StatusUnauthorized
	//log.Printf("%v", appUserIn.Username)
	roles := []string{"GUEST"}
	var appUser aumodel.AppUser
	if err := database.DB.Where("username = ?", appUserIn.Username).First(&appUser); err.Error != nil {
		log.Printf("User not found: %v error: %v\n", appUserIn.Username, err.Error)
		return result, status
	}
	if err := bcrypt.CompareHashAndPassword([]byte(appUser.Password), []byte(appUserIn.Password)); err != nil {
		log.Printf("Password wrong. Username: %v\n", appUser.Username)
		return result, status
	}
	// add jwt token creation
	var myUuid uuid.UUID
	if myUuid1, err := uuid.NewRandom(); err != nil {
		log.Printf("Uuid creation failed: %v\n", err.Error())
		return result, status
	} else {
		myUuid = myUuid1
		roles = append(roles, "USER")
	}

	tokenTtl := 60
	if len(strings.TrimSpace(os.Getenv("MSG_MESSAGES"))) <= 3 {
		tokenTtl = tokenTtl * 10
	}
	myToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":     appUser.Username,
		"uuid":    myUuid.String(),
		"auth":    strings.Join(roles[:], ","),
		"lastmsg": time.Now().Unix(),
		"exp":     time.Now().Add(time.Second * time.Duration(tokenTtl)).Unix(),
	})
	jwtTokenSecrect := os.Getenv("JWT_TOKEN_SECRET")
	result, err := myToken.SignedString([]byte(jwtTokenSecrect))
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

func generatePasswordHash(password string) []byte {
	passwordSlice := []byte(password)
	hashValue, err := bcrypt.GenerateFromPassword(passwordSlice, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("GenerateFromPassword failed: %v\n", err.Error())
	}
	return hashValue
}
