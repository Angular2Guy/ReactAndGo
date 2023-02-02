package gsmodel

import (
	"time"
)

type GasStation struct {
	ID                      string `gorm:"primaryKey"`
	Version                 string
	VersionTime             time.Time
	StationName             string `gorm:"column:name"`
	Brand                   string `gorm:"index:idx_brand"`
	Street                  string
	Place                   string
	HouseNumber             string
	PostCode                string  `gorm:"index:idx_gas_station_post_code"`
	Latitude                float64 `gorm:"column:lat;index:idx_lat"`
	Longitude               float64 `gorm:"column:lng;index:idx_lng"`
	PublicHolidayIdentifier string
	PriceInImport           time.Time `gorm:"index:idx_updated"`
	PriceChanged            time.Time
	OpenTs                  int `gorm:"index:idx_open_ts"`
	OtJson                  string
	StationInImport         time.Time
	FirstActive             time.Time
	GasPrices               []GasPrice
}
