package controller

import (
	"angular-and-go/pkd/appuser"
	aufile "angular-and-go/pkd/appuser/file"
	aubody "angular-and-go/pkd/controller/aumodel"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPlzCoordinates(c *gin.Context) {
	filePath := c.Params.ByName("path")
	aufile.UpdatePlzCoordinates(filePath)

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
