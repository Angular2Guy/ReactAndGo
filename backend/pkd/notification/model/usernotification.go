package unmodel

import "time"

type UserNotification struct {
	ID        int64     `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"index:idx_un_timestamp"`
	UserUuid  string    `gorm:"size:64;not null;index:idx_un_user_uuid"`
	Title     string    `gorm:"size:256"`
	Message   string    `gorm:"size:1024"`
	DataJson  string    `gorm:"size:4096"`
}
