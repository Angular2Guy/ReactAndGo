package controller

import (
	"angular-and-go/pkd/appuser"
	"angular-and-go/pkd/appuser/aumodel"
	aufile "angular-and-go/pkd/appuser/file"
	aubody "angular-and-go/pkd/controller/aumodel"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type postCodeLocationJson struct {
	Longitude  float64
	Latitude   float64
	Label      string
	PostCode   int32
	SquareKM   float32
	Population int32
}

func getLocation(c *gin.Context) {
	locationStr := c.Query("location")
	postCodeLocations := appuser.FindLocation(locationStr)
	myPostCodeLocations := mapToPostCodeLocation(postCodeLocations)
	c.JSON(http.StatusOK, myPostCodeLocations)
}

func mapToPostCodeLocation(postCodeLocations []aumodel.PostCodeLocation) []postCodeLocationJson {
	result := []postCodeLocationJson{}
	for _, postCodeLocation := range postCodeLocations {
		myPostCodeLocation := postCodeLocationJson{
			Longitude:  postCodeLocation.CenterLongitude,
			Latitude:   postCodeLocation.CenterLatitude,
			Label:      postCodeLocation.Label,
			PostCode:   postCodeLocation.PostCode,
			SquareKM:   postCodeLocation.SquareKM,
			Population: postCodeLocation.Population,
		}
		result = append(result, myPostCodeLocation)
	}
	return result
}

func getPostCodeCoordinates(c *gin.Context) {
	filePath := c.Query("filename")
	aufile.UpdatePostCodeCoordinates(filePath)
}

func postSignin(c *gin.Context) {
	//jsonData, err := ioutil.ReadAll(c.Request.Body)
	//fmt.Printf("Json: %v, Err: %v", string(jsonData), err)
	var appUserRequest aubody.AppUserRequest
	if err := c.Bind(&appUserRequest); err != nil {
		log.Printf("postSingin: %v", err.Error())
	}
	myAppUser := appuser.AppUserIn{Username: appUserRequest.Username, Password: appUserRequest.Password, Latitude: appUserRequest.Latitude, Uuid: "", Longitude: appUserRequest.Longitude}
	result := appuser.Signin(myAppUser)
	httpResult := http.StatusNotAcceptable
	message := ""
	if result == appuser.Ok {
		httpResult = http.StatusAccepted
	} else if result == appuser.UsernameTaken {
		message = "Username not available."
	}
	c.JSON(httpResult, aubody.AppUserResponse{Token: "", Message: message})
}

func postLogin(c *gin.Context) {
	var appUserRequest aubody.AppUserRequest
	if err := c.Bind(&appUserRequest); err != nil {
		log.Printf("postLogin: %v", err.Error())
	}
	myAppUser := appuser.AppUserIn{Username: appUserRequest.Username, Password: appUserRequest.Password, Latitude: appUserRequest.Latitude, Uuid: "", Longitude: appUserRequest.Longitude}
	result, status := appuser.Login(myAppUser)
	var message = ""
	if status != http.StatusOK {
		message = "Login failed."
	}
	appAuResponse := aubody.AppUserResponse{Token: result, Message: message}
	c.JSON(status, appAuResponse)
}
