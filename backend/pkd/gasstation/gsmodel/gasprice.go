package gsmodel

import "time"

type Tabler interface {
	TableName() string
}

type GasPrice struct {
	ID           int64  `gorm:"primaryKey"`
	GasStationID string `gorm:"column:stid"`
	E5           int
	E10          int
	Diesel       int
	Date         time.Time
	Changed      int
}

func (GasPrice) TableName() string {
	return "gas_station_information_history"
}
