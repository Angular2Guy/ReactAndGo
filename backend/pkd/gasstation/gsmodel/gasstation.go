package gsmodel

import (
	"time"
)

type GasStation struct {
	ID                      string `gorm:"primaryKey"`
	Version                 int
	VersionTime             time.Time
	StationName             string `gorm:"name"`
	Brand                   string
	Street                  string
	Place                   string
	HouseNumber             string
	PostCode                string
	Latitude                float64
	Longitude               float64
	PublicHolydayIdentifier string
	PriceInImport           time.Time
	PriceChanged            time.Time
	OpenTs                  int
	OtJson                  string
	StationInImport         time.Time
	FirstActive             time.Time
}
