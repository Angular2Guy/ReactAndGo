/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package gasstation

import (
	"fmt"
	"log"
	"os"
	gsbody "react-and-go/pkd/controller/gsmodel"
	"react-and-go/pkd/database"
	"react-and-go/pkd/gasstation/gsmodel"
	"react-and-go/pkd/notification"
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

func UpdateGasStations(gasStations *[]GasStationImport) {
	gasStationImportMap := make(map[string]GasStationImport)
	for _, value := range *gasStations {
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
			tx.Save(&resultGs)
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

func UpdatePrice(gasStationPrices *[]GasStationPrices) {
	stationPricesMap := make(map[string]GasStationPrices)
	var stationPricesKeys []string
	for _, value := range *gasStationPrices {
		stationPricesMap[value.GasStationID] = value
		stationPricesKeys = append(stationPricesKeys, value.GasStationID)
	}
	gasPriceUpdateMap := make(map[string]gsmodel.GasPrice)
	stationPricesDb := FindPricesByStids(&stationPricesKeys, 0)
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
			if stationPricesMap[value.GasStationID].Timestamp.Before(time.Now().Add(time.Hour*-720)) || stationPricesMap[value.GasStationID].Diesel < 0 ||
				stationPricesMap[value.GasStationID].E10 < 0 || stationPricesMap[value.GasStationID].E5 < 0 {
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
	if len(stationPricesMap) > 0 {
		var stationIds []string
		for _, stationPrice := range stationPricesMap {
			stationIds = append(stationIds, stationPrice.GasStationID)
		}
		myGasStations := findByIds(&stationIds)
		for _, gasStation := range myGasStations {
			value := stationPricesMap[gasStation.ID]
			gasPriceUpdateMap[value.GasStationID] = gsmodel.GasPrice{GasStationID: value.GasStationID, E5: stationPricesMap[value.GasStationID].E5, E10: stationPricesMap[value.GasStationID].E10,
				Diesel: stationPricesMap[value.GasStationID].Diesel, Date: stationPricesMap[value.GasStationID].Timestamp, Changed: 21}
			log.Printf("GasStation with first price: %v\n", gasStation.ID)
			delete(stationPricesMap, value.GasStationID)
		}
		//create new gas stations
		if len(stationPricesMap) > 0 {
			log.Default().Printf("New GasStations: %v\n", len(stationPricesMap))
		}
	}
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, value := range gasPriceUpdateMap {
			tx.Save(&value)
		}
		return nil
	})
	log.Printf("Prices updated: %v\n", len(gasPriceUpdateMap))
	go sendNotifications(&gasPriceUpdateMap)
	go updateCountyStatePrices(&gasPriceUpdateMap)
}

func updateCountyStatePrices(gasStationIDToGasPriceMap *map[string]gsmodel.GasPrice) int {
	var gasStationIDs []string
	for gasStationID, _ := range *gasStationIDToGasPriceMap {
		gasStationIDs = append(gasStationIDs, gasStationID)
	}
	gasStationIDChunks := createChunks(&gasStationIDs)
	var gasStations []gsmodel.GasStation
	for gasStationIDChunk := range gasStationIDChunks {
		var values []gsmodel.GasStation
		database.DB.Where("ID IN ?", gasStationIDChunk).Find(&values)
		gasStations = append(gasStations, values...)
	}
	return len(gasStations)
}

func sendNotifications(gasStationIDToGasPriceMap *map[string]gsmodel.GasPrice) {
	var gasStationIds []string
	for key, _ := range *gasStationIDToGasPriceMap {
		gasStationIds = append(gasStationIds, key)
	}
	gasStations := findByIds(&gasStationIds)
	notification.SendNotifications(gasStationIDToGasPriceMap, gasStations)
}

func createChunks(ids *[]string) [][]string {
	cunckedSelects := strings.ToLower(strings.TrimSpace(os.Getenv("DB_CHUNKED_SELECTS")))
	chunkSize := 10000
	if cunckedSelects == "true" {
		chunkSize = 999
	}
	chuncks := chunkSlice(*ids, chunkSize)
	if len(chuncks) > 1 {
		log.Printf("Number of Chunks: %v\n", len(chuncks))
	}
	return chuncks
}

func findByIds(ids *[]string) []gsmodel.GasStation {
	var result []gsmodel.GasStation
	chuncks := createChunks(ids)
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, chunk := range chuncks {
			var values []gsmodel.GasStation
			tx.Where("id in ?", chunk).Find(&values)
			result = append(result, values...)
		}
		return nil
	})
	return result
}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).First(&myGasStation)
	return myGasStation
}

func FindPricesByStids(stids *[]string, resultLimit int) []gsmodel.GasPrice {
	var myGasPrice []gsmodel.GasPrice
	oneMonthAgo := time.Now().Add(time.Hour * -720)
	dateStr := fmt.Sprintf("%04d-%02d-%02d", oneMonthAgo.Year(), oneMonthAgo.Month(), oneMonthAgo.Day())
	//log.Printf("Cut off date: %v", dateStr)
	chuncks := createChunks(stids)
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, chunk := range chuncks {
			var values []gsmodel.GasPrice
			//log.Printf("Chunk: %v\n", chunk)
			myQuery := tx.Where("stid IN ? and date >= date(?) ", chunk, dateStr).Order("date desc")
			if resultLimit > 0 {
				myQuery.Limit(resultLimit)
			}
			myQuery.Find(&values)
			myGasPrice = append(myGasPrice, values...)
		}
		return nil
	})
	return myGasPrice
}

func FindByPostCodes(postcodes []string) []gsmodel.GasStation {
	return nil
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
	query.Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(50)
	}).Find(&gasStations)
	return gasStations
}

func FindBySearchLocation(searchLocation gsbody.SearchLocation) []gsmodel.GasStation {
	var gasStations []gsmodel.GasStation
	myRadius := searchLocation.Radius + 0.1
	if myRadius > 20.0 {
		myRadius = 20.1
	}
	minMax := calcMinMaxSquare(searchLocation.Longitude, searchLocation.Latitude, myRadius)
	//fmt.Printf("WestLat: %v, WestLng: %v\n", westLat, westLng)
	//fmt.Printf("MinLat: %v, MinLng: %v, MaxLat: %v, MaxLng: %v\n", minMax.MinLat, minMax.MinLng, minMax.MaxLat, minMax.MaxLng)
	database.DB.Where("lat >= ? and lat <= ? and lng >= ? and lng <= ?", minMax.MinLat, minMax.MaxLat, minMax.MinLng, minMax.MaxLng).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(50)
	}).Find(&gasStations)
	//filter for stations in circle
	filteredGasStations := []gsmodel.GasStation{}
	for _, myGasStation := range gasStations {
		distance, bearing := myGasStation.CalcDistanceBearing(searchLocation.Latitude, searchLocation.Longitude)
		//fmt.Printf("Distance: %v, Bearing: %v\n", distance, bearing)
		if distance < myRadius && bearing > -1.0 {
			filteredGasStations = append(filteredGasStations, myGasStation)
		}
	}
	return filteredGasStations
}

func calcMinMaxSquare(longitude float64, latitude float64, radius float64) minMaxSquare {
	minMax := minMaxSquare{MinLat: 1000.0, MinLng: 1000.0, MaxLat: 0.0, MaxLng: 0.0}
	//fmt.Printf("StartLat: %v, StartLng: %v Radius: %v\n", searchLocation.Latitude, searchLocation.Longitude, searchLocation.Radius)
	//max supported radius 20km and add 0.1 for floation point side effects
	northLat, northLng := gsmodel.CalcLocation(latitude, longitude, radius, 0.0)
	minMax = updateMinMaxSquare(northLat, northLng, minMax)
	//fmt.Printf("NorthLat: %v, NorthLng: %v\n", northLat, northLng)
	eastLat, eastLng := gsmodel.CalcLocation(latitude, longitude, radius, 90.0)
	minMax = updateMinMaxSquare(eastLat, eastLng, minMax)
	//fmt.Printf("EastLat: %v, EastLng: %v\n", eastLat, eastLng)
	southLat, southLng := gsmodel.CalcLocation(latitude, longitude, radius, 180.0)
	minMax = updateMinMaxSquare(southLat, southLng, minMax)
	//fmt.Printf("SouthLat: %v, SouthLng: %v\n", southLat, southLng)
	westLat, westLng := gsmodel.CalcLocation(latitude, longitude, radius, 270.0)
	minMax = updateMinMaxSquare(westLat, westLng, minMax)
	return minMax
}

func chunkSlice[T any](mySlice []T, chunkSize int) (s [][]T) {
	numberOfChunks := len(mySlice)/chunkSize + 1
	var result [][]T
	for i := 0; i < numberOfChunks; i++ {
		min := (i * len(mySlice) / numberOfChunks)
		max := ((i + 1) * len(mySlice)) / numberOfChunks
		result = append(result, mySlice[min:max])
	}
	return result
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
