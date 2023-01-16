package aumodel

import "time"

type LoggedOutUser struct {
	ID         int64  `gorm:"primaryKey"`
	Username   string `gorm:"size:64;not null;index:idx_lou_user_name"`
	Uuid       string `gorm:"size:64;not null;index:idx_lou_uuid,unique"`
	LastLogout time.Time
}
