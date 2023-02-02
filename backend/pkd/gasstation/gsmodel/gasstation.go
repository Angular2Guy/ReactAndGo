package gsmodel

import (
	"math"
	"time"
)

const earthRadius = 6371.0

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

type MyGasStation interface {
	GasStation
	CalcDistance(startLat float64, startLng float64) (float64, float64)
}

func (gasStation GasStation) CalcDistanceBearing(startLat float64, startLng float64) (float64, float64) {
	var radStartLat = toRad1(startLat)
	var radDestLat = toRad1(gasStation.Latitude)
	var radDeltaLat = toRad1(gasStation.Latitude - startLat)
	var radDeltaLng = toRad1(gasStation.Longitude - startLng)
	//distance
	var a = math.Sin(radDeltaLat/2)*math.Sin(radDeltaLat/2) + math.Cos(radStartLat)*math.Cos(radDestLat)*math.Sin(radDeltaLng/2)*math.Sin(radDeltaLng/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var distance = earthRadius * c
	//bearing
	var y = math.Sin(radDeltaLng) * math.Cos(radDestLat)
	var x = math.Cos(radStartLat)*math.Sin(radDestLat) - math.Sin(radStartLat)*math.Cos(radDestLat)*math.Cos(radDeltaLng)
	var bearing = math.Mod((toDeg1(math.Atan2(y, x)) + 360.0), 360.0)
	return distance, bearing
}

func toRad1(myValue float64) float64 {
	return myValue * math.Pi / 180
}

func toDeg1(myValue float64) float64 {
	return myValue * 180 / math.Pi
}
