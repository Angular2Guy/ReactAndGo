package gasstation

import (
	gsbody "angular-and-go/pkd/contr/model"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/gasstation/gsmodel"
	"math"
	"strings"

	"gorm.io/gorm"
)

const earthRadius = 6371.0

type minMaxSquare struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).First(&myGasStation)
	return myGasStation
}

func FindPricesByStid(stid string) []gsmodel.GasPrice {
	var myGasPrice []gsmodel.GasPrice
	database.DB.Where("stid = ?", stid).Order("date desc").Find(&myGasPrice)
	return myGasPrice
}

func FindBySearchPlace(searchPlace gsbody.SearchPlaceBody) []gsmodel.GasStation {
	var gasStations []gsmodel.GasStation
	var query = database.DB
	if len(strings.TrimSpace(searchPlace.Place)) >= 2 {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(searchPlace.Place)+"%")
	}
	if len(strings.TrimSpace(searchPlace.PostCode)) >= 4 {
		query = query.Where("post_code LIKE ?", "%"+strings.TrimSpace(searchPlace.PostCode)+"%")
	}
	if len(strings.TrimSpace(searchPlace.StationName)) >= 2 {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(searchPlace.StationName)+"%")
	}
	query.Limit(200).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).Find(&gasStations)
	return gasStations
}

func FindBySearchLocation(searchLocation gsbody.SearchLocation) []gsmodel.GasStation {
	var gasStations []gsmodel.GasStation
	minMax := minMaxSquare{MinLat: 1000.0, MinLng: 1000.0, MaxLat: 0.0, MaxLng: 0.0}
	//fmt.Printf("StartLat: %v, StartLng: %v\n", searchLocation.Latitude, searchLocation.Longitude)
	northLat, northLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, 20.0, 0.0)
	minMax = updateMinMaxSquare(northLat, northLng, minMax)
	//fmt.Printf("NorthLat: %v, NorthLng: %v\n", northLat, northLng)
	eastLat, eastLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, 20.0, 90.0)
	minMax = updateMinMaxSquare(eastLat, eastLng, minMax)
	//fmt.Printf("EastLat: %v, EastLng: %v\n", eastLat, eastLng)
	southLat, southLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, 20.0, 180.0)
	minMax = updateMinMaxSquare(southLat, southLng, minMax)
	//fmt.Printf("SouthLat: %v, SouthLng: %v\n", southLat, southLng)
	westLat, westLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, 20.0, 270.0)
	minMax = updateMinMaxSquare(westLat, westLng, minMax)
	//fmt.Printf("WestLat: %v, WestLng: %v\n", westLat, westLng)
	//fmt.Printf("MinLat: %v, MinLng: %v, MaxLat: %v, MaxLng: %v\n", minMax.MinLat, minMax.MinLng, minMax.MaxLat, minMax.MaxLng)
	database.DB.Where("lat >= ? and lat <= ? and lng >= ? and lng <= ?", minMax.MinLat, minMax.MaxLat, minMax.MinLng, minMax.MaxLng).Limit(200).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).Find(&gasStations)
	//filter for stations in circle
	filteredGasStations := []gsmodel.GasStation{}
	for _, myGasStation := range gasStations {
		distance, bearing := calcDistance(searchLocation.Latitude, searchLocation.Longitude, myGasStation.Latitude, myGasStation.Longitude)
		//fmt.Printf("Distance: %v, Bearing: %v\n", distance, bearing)
		if distance < 20.1 && bearing > -1.0 {
			filteredGasStations = append(filteredGasStations, myGasStation)
		}
	}
	return filteredGasStations
}

func calcDistance(startLat float64, startLng float64, destLat float64, destLng float64) (float64, float64) {
	var radStartLat = toRad(startLat)
	var radDestLat = toRad(destLat)
	var radDeltaLat = toRad(destLat - startLat)
	var radDeltaLng = toRad(destLng - startLng)
	//distance
	var a = math.Sin(radDeltaLat/2)*math.Sin(radDeltaLat/2) + math.Cos(radStartLat)*math.Cos(radDestLat)*math.Sin(radDeltaLng/2)*math.Sin(radDeltaLng/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var distance = earthRadius * c
	//bearing
	var y = math.Sin(radDeltaLng) * math.Cos(radDestLat)
	var x = math.Cos(radStartLat)*math.Sin(radDestLat) - math.Sin(radStartLat)*math.Cos(radDestLat)*math.Cos(radDeltaLng)
	var bearing = math.Mod((toDeg(math.Atan2(y, x)) + 360.0), 360.0)
	return distance, bearing
}

func calcLocation(startLat float64, startLng float64, distanceKm float64, bearing float64) (float64, float64) {
	var radBearing = toRad(bearing)
	var radStartLat = toRad(startLat)
	var radStartLng = toRad(startLng)
	var radDestLat = math.Asin(math.Sin(radStartLat)*math.Cos(distanceKm/earthRadius) + math.Cos(radStartLat)*math.Sin(distanceKm/earthRadius)*math.Cos(radBearing))
	var radDestLng = radStartLng + math.Atan2(math.Sin(radBearing)*math.Sin(distanceKm/earthRadius)*math.Cos(radStartLat), math.Cos(distanceKm/earthRadius)-math.Sin(radStartLat)*math.Sin(radDestLat))
	destLat := toDeg(radDestLat)
	destLng := toDeg(radDestLng)
	return destLat, destLng
}

func toRad(myValue float64) float64 {
	return myValue * math.Pi / 180
}

func toDeg(myValue float64) float64 {
	return myValue * 180 / math.Pi
}

func updateMinMaxSquare(newLat float64, newLng float64, minMax minMaxSquare) minMaxSquare {
	if newLat > minMax.MaxLat {
		minMax.MaxLat = newLat
	}
	if newLat < minMax.MinLat {
		minMax.MinLat = newLat
	}
	if newLng > minMax.MaxLng {
		minMax.MaxLng = newLng
	}
	if newLng < minMax.MinLng {
		minMax.MinLng = newLng
	}
	return minMax
}
