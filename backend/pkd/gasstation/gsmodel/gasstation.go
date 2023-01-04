package gsmodel

import (
	"time"
)

type GasStation struct {
	ID                      string `gorm:"primaryKey"`
	Version                 int
	VersionTime             time.Time
	StationName             string `gorm:"column:name"`
	Brand                   string
	Street                  string
	Place                   string
	HouseNumber             string
	PostCode                string
	Latitude                float64 `gorm:"column:lat"`
	Longitude               float64 `gorm:"column:lng"`
	PublicHolydayIdentifier string
	PriceInImport           time.Time
	PriceChanged            time.Time
	OpenTs                  int
	OtJson                  string
	StationInImport         time.Time
	FirstActive             time.Time
	//GasPrices               []GasPrice
}
