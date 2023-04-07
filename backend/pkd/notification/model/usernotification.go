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
package unmodel

import (
	"time"
)

type UserNotification struct {
	ID               int64     `gorm:"primaryKey"`
	Timestamp        time.Time `gorm:"index:idx_un_timestamp"`
	UserUuid         string    `gorm:"size:64;not null;index:idx_un_user_uuid"`
	Title            string    `gorm:"size:256"`
	Message          string    `gorm:"size:4096"`
	DataJson         string
	NotificationSend bool
}
