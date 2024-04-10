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
package controller

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"react-and-go/pkd/appuser"
	aubody "react-and-go/pkd/controller/aumodel"
	fileim "react-and-go/pkd/fileimport"
	postcode "react-and-go/pkd/postcode"
	pcmodel "react-and-go/pkd/postcode/pcmodel"
	token "react-and-go/pkd/token"

	"github.com/gin-gonic/gin"
)

func getLogout(c *gin.Context) {
	status := http.StatusUnauthorized
	message := "Invalid"
	username, exists1 := c.Get("user")
	uuid, exists2 := c.Get("uuid")
	if exists1 && exists2 {
		token.LoggedOutUsers = appuser.StoreUserLogout(username.(string), uuid.(string))
		if len(token.LoggedOutUsers) > 0 {
			message = ""
			status = http.StatusOK
		}
	}
	c.JSON(status, aubody.AppUserResponse{Token: "", Message: message})
}

func getRefreshToken(c *gin.Context) {
	status := http.StatusUnauthorized
	message := "Invalid"
	result := ""
	userName, exits := c.Get("user")
	roles, exists2 := c.Get("roles")
	if exits && exists2 {
		//jwt token creation
		var err error
		result, err = token.CreateToken(token.TokenUser{Username: userName.(string), Roles: []string{roles.(string)}})
		if err != nil {
			log.Printf("Failed to create jwt token: %v\n", err)
		} else {
			status = http.StatusOK
			message = ""
		}
	}
	c.JSON(status, aubody.AppUserResponse{Token: result, Message: message})
}

func getLocation(c *gin.Context) {
	locationStr := c.Query("location")
	postCodeLocations := postcode.FindLocation(locationStr)
	//log.Printf("Locations: %v", postCodeLocations)
	myPostCodeLocations := mapToPostCodeLocation(postCodeLocations)
	c.JSON(http.StatusOK, myPostCodeLocations)
}

func mapToPostCodeLocation(postCodeLocations []pcmodel.PostCodeLocation) []aubody.CodeLocationResponse {
	result := []aubody.CodeLocationResponse{}
	for _, postCodeLocation := range postCodeLocations {
		if !math.IsNaN(postCodeLocation.CenterLatitude) && !math.IsNaN(postCodeLocation.CenterLongitude) && !math.IsNaN(float64(postCodeLocation.SquareKM)) {
			myPostCodeLocation := aubody.CodeLocationResponse{
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
	fileim.UpdatePostCodeCoordinates(filePath)
}

func getStateCountyData(c *gin.Context) {
	filePath := c.Query("filename")
	fileim.UpdateStatesAndCounties(filePath)
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
	result, status, postCode, userUuid, userLongitude, userLatitude, searchRadius, targetE5, targetE10, targetDiesel := appuser.Login(myAppUser)
	var message = ""
	if status != http.StatusOK {
		message = "Login failed."
	}
	appAuResponse := aubody.AppUserResponse{Token: result, Message: message, PostCode: postCode, Uuid: userUuid, Longitude: userLongitude, Latitude: userLatitude,
		SearchRadius: searchRadius, TargetE5: fmt.Sprintf("%v", (float64(targetE5) / 1000)), TargetE10: fmt.Sprintf("%v", (float64(targetE10) / 1000)), TargetDiesel: fmt.Sprintf("%v", (float64(targetDiesel) / 1000))}
	c.JSON(status, appAuResponse)
}

func postUserLocationRadius(c *gin.Context) {
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
	c.JSON(httpResult, aubody.CodeLocationResponse{Message: message, Label: "", Longitude: appUserRequest.Longitude, Latitude: appUserRequest.Latitude, PostCode: 0, SquareKM: 0, Population: 0})
}

func postTargetPrices(c *gin.Context) {
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
	c.JSON(httpResult, aubody.TargetPricesResponse{Message: message, TargetDiesel: appUserRequest.TargetDiesel, TargetE10: appUserRequest.TargetE10, TargetE5: appUserRequest.TargetE5})
}
