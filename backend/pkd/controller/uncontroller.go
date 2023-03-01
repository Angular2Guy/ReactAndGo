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
