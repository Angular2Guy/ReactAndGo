package gsmodel

import (
	"gorm.io/gorm"
)

type GasStation struct {
	gorm.Model
	StationName string
}
