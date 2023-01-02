package gsmodel

import "time"

type Tabler interface {
	TableName() string
}

type GasPrice struct {
	ID      int64 `gorm:"primaryKey"`
	Stid    string
	E5      int
	E10     int
	Diesel  int
	date    time.Time
	changed int
}

func (GasPrice) TableName() string {
	return "gas_station_information_history"
}
