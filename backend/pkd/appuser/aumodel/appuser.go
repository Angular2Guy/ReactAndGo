package aumodel

import (
	"gorm.io/gorm"
)

type AppUser struct {
	gorm.Model
	Username     string `gorm:"size:64;not null;index:idx_au_user_name,unique"`
	Password     string `gorm:"size:32;not null"`
	Latitude     float64
	Longitude    float64
	SearchRadius float64
	TargetDiesel int
	TargetE5     int
	TargetE10    int
}
