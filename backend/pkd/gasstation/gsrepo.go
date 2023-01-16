package gasstation

import (
	gsbody "angular-and-go/pkd/contr/gsmodel"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/gasstation/gsmodel"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

const earthRadius = 6371.0

type minMaxSquare struct {
	MinLat float64
	MinLng float64
	MaxLat float64
	MaxLng float64
}

type GasStationPrices struct {
	GasStationID string `gorm:"column:stid"`
	E5           int
	E10          int
	Diesel       int
	Timestamp    time.Time
}

type GasStationImport struct {
	Uuid             string
	StationName      string
	Brand            string
	Street           string
	HouseNumber      string
	PostCode         string
	City             string
	Latitude         float64
	Longitude        float64
	FirstActive      time.Time
	OpeningTimesJson string
}

func UpdateGasStations(gasStations []GasStationImport) {
	gasStationImportMap := make(map[string]GasStationImport)
	for _, value := range gasStations {
		gasStationImportMap[value.Uuid] = value
	}
	fmt.Printf("GasStations found: %v\n", len(gasStationImportMap))
	gasStationsUpdated := 0
	var values []gsmodel.GasStation
	database.DB.Transaction(func(tx *gorm.DB) error {
		tx.FindInBatches(&values, 1000, func(tx *gorm.DB, batch int) error {
			var newResults []gsmodel.GasStation
			for _, result := range values {
				if strings.TrimSpace(result.OtJson) != strings.TrimSpace(gasStationImportMap[result.ID].OpeningTimesJson) {
					result.OtJson = gasStationImportMap[result.ID].OpeningTimesJson
					newResults = append(newResults, result)
				}
				delete(gasStationImportMap, result.ID)
			}
			if len(newResults) > 0 {
				tx.Save(&newResults)
			}
			gasStationsUpdated = gasStationsUpdated + len(newResults)
			//tx.RowsAffected // number of records in this batch

			//batch // Batch 1, 2, 3

			// returns error will stop future batches
			return nil
		})
		fmt.Printf("GasStations updated: %v\n", gasStationsUpdated)
		for _, value := range gasStationImportMap {
			resultGs := createNewGasStation(value)
			tx.Save(resultGs)
		}
		fmt.Printf("GasStations new: %v\n", len(gasStationImportMap))
		return nil
	})
}

func createNewGasStation(value GasStationImport) gsmodel.GasStation {
	var resultGs gsmodel.GasStation
	resultGs.ID = value.Uuid
	resultGs.Brand = value.Brand
	resultGs.FirstActive = value.FirstActive
	resultGs.HouseNumber = value.HouseNumber
	resultGs.Latitude = value.Latitude
	resultGs.Longitude = value.Longitude
	resultGs.OpenTs = 0
	resultGs.OtJson = value.OpeningTimesJson
	resultGs.Place = value.City
	resultGs.PostCode = value.PostCode
	resultGs.PriceChanged = time.Now()
	resultGs.PriceInImport = time.Now()
	//resultGs.PublicHolidayIdentifier = ""
	resultGs.StationInImport = time.Now()
	resultGs.StationName = value.StationName
	resultGs.Street = value.Street
	resultGs.Version = "1"
	resultGs.VersionTime = time.Now()
	return resultGs
}

func UpdatePrice(gasStationPrices []GasStationPrices) {
	stationPricesMap := make(map[string]GasStationPrices)
	var stationPricesKeys []string
	for _, value := range gasStationPrices {
		stationPricesMap[value.GasStationID] = value
		stationPricesKeys = append(stationPricesKeys, value.GasStationID)
	}
	gasPriceUpdateMap := make(map[string]gsmodel.GasPrice)
	stationPricesDb := FindPricesByStids(stationPricesKeys)
	log.Printf("StationPricesKeys: %v StationPricesDb: %v", len(stationPricesKeys), len(stationPricesDb))
	for _, value := range stationPricesDb {
		if _, found := gasPriceUpdateMap[value.GasStationID]; !found {
			var myChanges = 0
			if stationPricesMap[value.GasStationID].Diesel != value.Diesel {
				myChanges = myChanges + 1
			}
			if stationPricesMap[value.GasStationID].E10 != value.E10 {
				myChanges = myChanges + 16
			}
			if stationPricesMap[value.GasStationID].E5 != value.E5 {
				myChanges = myChanges + 4
			}
			// validation checks
			if stationPricesMap[value.GasStationID].Timestamp.Before(time.Now().Add(time.Hour*-720)) || stationPricesMap[value.GasStationID].Diesel < 15 ||
				stationPricesMap[value.GasStationID].E10 < 15 || stationPricesMap[value.GasStationID].E5 < 15 {
				myChanges = 0
			}
			//log.Printf("GasStation: %v Changes: %v", value.GasStationID, myChanges)
			if myChanges > 0 {
				gasPriceUpdateMap[value.GasStationID] = gsmodel.GasPrice{GasStationID: value.GasStationID, E5: stationPricesMap[value.GasStationID].E5, E10: stationPricesMap[value.GasStationID].E10,
					Diesel: stationPricesMap[value.GasStationID].Diesel, Date: stationPricesMap[value.GasStationID].Timestamp, Changed: myChanges}
				//value, _ := json.Marshal(gasPriceUpdateMap[value.GasStationID])
				//fmt.Printf("Update: %v\n", string(value))
			}
			delete(stationPricesMap, value.GasStationID)
		}
	}
	for _, value := range gasPriceUpdateMap {
		database.DB.Save(&value)
	}
	log.Printf("Prices updated: %v\n", len(gasPriceUpdateMap))
	if len(stationPricesMap) > 0 {
		//create new gas stations
		log.Default().Printf("New GasStations: %v\n", len(stationPricesMap))
	}

}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).First(&myGasStation)
	return myGasStation
}

func FindPricesByStids(stids []string) []gsmodel.GasPrice {
	var myGasPrice []gsmodel.GasPrice
	threeMonthsAgo := time.Now().Add(time.Hour * -2160)
	dateStr := fmt.Sprintf("%04d-%02d-%02d", threeMonthsAgo.Year(), threeMonthsAgo.Month(), threeMonthsAgo.Day())
	//log.Printf("Cut off date: %v", dateStr)
	database.DB.Where("stid IN ? and date >= date(?) ", stids, dateStr).Order("date desc").Find(&myGasPrice)
	return myGasPrice
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
	//fmt.Printf("StartLat: %v, StartLng: %v Radius: %v\n", searchLocation.Latitude, searchLocation.Longitude, searchLocation.Radius)
	//max supported radius 20km and add 0.1 for floation point side effects
	myRadius := searchLocation.Radius + 0.1
	if myRadius > 20.0 {
		myRadius = 20.1
	}
	northLat, northLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, myRadius, 0.0)
	minMax = updateMinMaxSquare(northLat, northLng, minMax)
	//fmt.Printf("NorthLat: %v, NorthLng: %v\n", northLat, northLng)
	eastLat, eastLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, myRadius, 90.0)
	minMax = updateMinMaxSquare(eastLat, eastLng, minMax)
	//fmt.Printf("EastLat: %v, EastLng: %v\n", eastLat, eastLng)
	southLat, southLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, myRadius, 180.0)
	minMax = updateMinMaxSquare(southLat, southLng, minMax)
	//fmt.Printf("SouthLat: %v, SouthLng: %v\n", southLat, southLng)
	westLat, westLng := calcLocation(searchLocation.Latitude, searchLocation.Longitude, myRadius, 270.0)
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
		if distance < myRadius && bearing > -1.0 {
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
