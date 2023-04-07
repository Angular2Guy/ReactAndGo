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

func LoadNotifications(userUuid string, newNotifications bool) []unmodel.UserNotification {
	var userNotifications []unmodel.UserNotification
	if newNotifications {
		database.DB.Transaction(func(tx *gorm.DB) error {
			tx.Where("user_uuid = ? and notification_send = ?", userUuid, !newNotifications).Order("timestamp desc").Find(&userNotifications)
			for _, userNotification := range userNotifications {
				userNotification.NotificationSend = true
				tx.Save(&userNotification)
			}
			return nil
		})
	} else {
		database.DB.Transaction(func(tx *gorm.DB) error {
			tx.Where("user_uuid = ?", userUuid).Order("timestamp desc").Find(&userNotifications)
			var myUserNotifications []unmodel.UserNotification
			for index, userNotification := range userNotifications {
				if index < 10 {
					myUserNotifications = append(myUserNotifications, userNotification)
					continue
				}
				tx.Delete(&userNotification)
			}
			userNotifications = myUserNotifications
			return nil
		})
	}
	return userNotifications
}
