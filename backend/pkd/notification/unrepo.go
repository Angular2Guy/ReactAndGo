package notification

import "log"

type NotificationMsg struct {
	UserUuid string
	Title    string
	Message  string
	DataJson string
}

func StoreNotifications(notificationMsgs []NotificationMsg) {
	for _, notificationMsg := range notificationMsgs {
		log.Printf("%v\n", notificationMsg.Title)
	}
}
