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
package aumodel

import "time"

type LoggedOutUser struct {
	ID         int64  `gorm:"primaryKey"`
	Username   string `gorm:"size:64;not null;index:idx_lou_user_name"`
	Uuid       string `gorm:"size:64;not null;index:idx_lou_uuid,unique"`
	LastLogout time.Time
}
