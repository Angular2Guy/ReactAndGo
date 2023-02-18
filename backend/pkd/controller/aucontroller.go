package controller

import (
	"angular-and-go/pkd/appuser"
	"angular-and-go/pkd/appuser/aumodel"
	aufile "angular-and-go/pkd/appuser/file"
	aubody "angular-and-go/pkd/controller/aumodel"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getLocation(c *gin.Context) {
	locationStr := c.Query("location")
	postCodeLocations := appuser.FindLocation(locationStr)
	//log.Printf("Locations: %v", postCodeLocations)
	myPostCodeLocations := mapToPostCodeLocation(postCodeLocations)
	c.JSON(http.StatusOK, myPostCodeLocations)
}

func mapToPostCodeLocation(postCodeLocations []aumodel.PostCodeLocation) []aubody.PostCodeLocationResponse {
	result := []aubody.PostCodeLocationResponse{}
	for _, postCodeLocation := range postCodeLocations {
		if !math.IsNaN(postCodeLocation.CenterLatitude) && !math.IsNaN(postCodeLocation.CenterLongitude) && !math.IsNaN(float64(postCodeLocation.SquareKM)) {
			myPostCodeLocation := aubody.PostCodeLocationResponse{
				Longitude:  postCodeLocation.CenterLongitude,
				Latitude:   postCodeLocation.CenterLatitude,
				Label:      postCodeLocation.Label,
				PostCode:   postCodeLocation.PostCode,
				SquareKM:   postCodeLocation.SquareKM,
				Population: postCodeLocation.Population,
			}
			result = append(result, myPostCodeLocation)
		}
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
	myAppUser := appuser.AppUserIn{Username: appUserRequest.Username, Password: appUserRequest.Password, Uuid: ""}
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
	myAppUser := appuser.AppUserIn{Username: appUserRequest.Username, Password: appUserRequest.Password, Uuid: ""}
	result, status, userLongitude, userLatitude, searchRadius, targetE5, targetE10, targetDiesel := appuser.Login(myAppUser)
	var message = ""
	if status != http.StatusOK {
		message = "Login failed."
	}
	appAuResponse := aubody.AppUserResponse{Token: result, Message: message, Longitude: userLongitude, Latitude: userLatitude,
		SearchRadius: searchRadius, TargetE5: fmt.Sprintf("%v", (targetE5 / 1000)), TargetE10: fmt.Sprintf("%v", (targetE10 / 1000)), TargetDiesel: fmt.Sprintf("%v", (targetDiesel / 1000))}
	c.JSON(status, appAuResponse)
}

func putUserLocationRadius(c *gin.Context) {
	var appUserRequest aubody.AppUserRequest
	if err := c.Bind(&appUserRequest); err != nil {
		log.Printf("putUserLocationRadius: %v", err.Error())
	}
	myAppUser := appuser.AppUserIn{Username: appUserRequest.Username, Uuid: "", Longitude: appUserRequest.Longitude, Latitude: appUserRequest.Latitude, SearchRadius: appUserRequest.SearchRadius}
	result := appuser.StoreLocationAndRadius(myAppUser)
	httpResult := http.StatusOK
	message := "Ok"
	if result != appuser.Ok {
		httpResult = http.StatusBadRequest
		message = "Invalid"
	}
	c.JSON(httpResult, aubody.AppUserResponse{Token: "", Message: message, Longitude: appUserRequest.Longitude, Latitude: appUserRequest.Latitude, SearchRadius: appUserRequest.SearchRadius, TargetDiesel: "0", TargetE10: "0", TargetE5: "0"})
}

func putTargetPrices(c *gin.Context) {
	var appUserRequest aubody.AppUserRequest
	if err := c.Bind(&appUserRequest); err != nil {
		log.Printf("putUserLocationRadius: %v", err.Error())
	}
	myTargetPrices := appuser.AppTargetIn{Username: appUserRequest.Username, TargetDiesel: appUserRequest.TargetDiesel, TargetE10: appUserRequest.TargetE10, TargetE5: appUserRequest.TargetE5}
	result := appuser.StoreTargetPrices(myTargetPrices)
	httpResult := http.StatusOK
	message := "Ok"
	if result != appuser.Ok {
		httpResult = http.StatusBadRequest
		message = "Invalid"
	}
	c.JSON(httpResult, aubody.AppUserResponse{Token: "", Message: message, Longitude: 0.0, Latitude: 0.0, SearchRadius: 0.0, TargetDiesel: appUserRequest.TargetDiesel, TargetE10: appUserRequest.TargetE10, TargetE5: appUserRequest.TargetE5})
}
