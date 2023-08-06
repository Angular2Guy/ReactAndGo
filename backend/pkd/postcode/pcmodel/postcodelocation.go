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
package pcmodel

import "gorm.io/gorm"

type PostCodeLocation struct {
	gorm.Model
	Label           string `gorm:"size:256;not null;index:idx_post_code_location_label"`
	PostCode        int32  `gorm:"index:idx_post_code_location_post_code"`
	Population      int32
	StateDataID     uint
	StateData       StateData
	CountyDataID    uint
	CountyData      CountyData
	SquareKM        float32
	CenterLongitude float64 `gorm:"index:idx_post_code_location_center_logitude"`
	CenterLatitude  float64 `gorm:"index:idx_post_code_location_center_latitude"`
}
