package notification

import (
	"log"
	"react-and-go/pkd/database"
	unmodel "react-and-go/pkd/notification/model"
	"time"

	"gorm.io/gorm"
)

type NotificationMsg struct {
	UserUuid string
	Title    string
	Message  string
	DataJson string
}

func StoreNotifications(notificationMsgs []NotificationMsg) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, notificationMsg := range notificationMsgs {
			log.Printf("%v\n", notificationMsg.Title)
			myUserNotification := unmodel.UserNotification{Timestamp: time.Now(), UserUuid: notificationMsg.UserUuid,
				Title: notificationMsg.Title, Message: notificationMsg.Message, DataJson: notificationMsg.DataJson, NotificationSend: false}
			tx.Save(&myUserNotification)
		}
		return nil
	})
}
