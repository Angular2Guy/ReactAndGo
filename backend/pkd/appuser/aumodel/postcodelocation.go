package aumodel

import "gorm.io/gorm"

type PostCodeLocation struct {
	gorm.Model
	Label           string `gorm:"size:256;not null;index:idx_post_code_location_label"`
	PostCode        int32  `gorm:"index:idx_post_code_location_post_code"`
	Population      int32
	SquareKM        float32
	CenterLongitude float64 `gorm:"index:idx_post_code_location_center_logitude"`
	CenterLatitude  float64 `gorm:"index:idx_post_code_location_center_latitude"`
}
