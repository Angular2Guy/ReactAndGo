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
	"net/http"
	unbody "react-and-go/pkd/controller/unmodel"
	notification "react-and-go/pkd/notification"
	unmodel "react-and-go/pkd/notification/model"

	"github.com/gin-gonic/gin"
)

func getNewUserNotifications(c *gin.Context) {
	userUuid := c.Param("useruuid")
	myNotifications := notification.LoadNotifications(userUuid, true)
	c.JSON(http.StatusOK, mapToUnResponses(myNotifications))
}

func getCurrentUserNotifications(c *gin.Context) {
	userUuid := c.Param("useruuid")
	myNotifications := notification.LoadNotifications(userUuid, false)
	c.JSON(http.StatusOK, mapToUnResponses(myNotifications))
}

func mapToUnResponses(myNotifications []unmodel.UserNotification) []unbody.UnResponse {
	var unResponses []unbody.UnResponse
	for _, myNotification := range myNotifications {
		unResponse := unbody.UnResponse{
			Timestamp: myNotification.Timestamp, UserUuid: myNotification.UserUuid, Title: myNotification.Title, Message: myNotification.Message, DataJson: myNotification.DataJson,
		}
		unResponses = append(unResponses, unResponse)
	}
	return unResponses
}
