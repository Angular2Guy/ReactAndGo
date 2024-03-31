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
	gsbody "react-and-go/pkd/controller/gsmodel"
	"react-and-go/pkd/database"
	"react-and-go/pkd/gasstation/gsmodel"
	"react-and-go/pkd/postcode"
	"react-and-go/pkd/postcode/pcmodel"
	"strings"
	"time"

	"gorm.io/gorm"
)

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

func UpdatePrice(gasStationPrices *[]GasStationPrices) {
	stationPricesMap := make(map[string]GasStationPrices)
	var stationPricesKeys []string
	for _, value := range *gasStationPrices {
		stationPricesMap[value.GasStationID] = value
		stationPricesKeys = append(stationPricesKeys, value.GasStationID)
	}
	gasPriceUpdateMap := make(map[string]gsmodel.GasPrice)
	stationPricesDb := FindPricesByStids(&stationPricesKeys, 0, Month)
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

func ReCalcCountyStatePrices() {
	log.Printf("recalcCountyStatePrices started.")
	myStart := time.Now()
	postCodePostCodeLocationMap, idStateDataMap, idCountyDataMap, postCodeGasStationsMap := createPostCodeGasStationMaps()
	gasStationIdGasPriceMap := createGasStationIdGasPriceMap(&postCodeGasStationsMap, Month)
	log.Printf("gasStationIdGasPriceMap: %v", len(gasStationIdGasPriceMap))
	resetDataMaps(&idStateDataMap, &idCountyDataMap)
	//sum up prices and count stations
	for _, myPostCodeLocation := range postCodePostCodeLocationMap {
		myPostCode := postcode.FormatPostCode(myPostCodeLocation.PostCode)
		for _, myGasStation := range postCodeGasStationsMap[myPostCode] {
			myGasPrice := gasStationIdGasPriceMap[myGasStation.ID]
			if myGasPrice.E5 < 10 && myGasPrice.E10 < 10 && myGasPrice.Diesel < 10 {
				continue
			}
			//log.Printf("%v", myGasPrice)
			myStateData := idStateDataMap[int(myPostCodeLocation.StateData.ID)]
			myCountyData := idCountyDataMap[int(myPostCodeLocation.CountyData.ID)]
			myStateData.GasStationNum += 1
			myCountyData.GasStationNum += 1
			if myGasPrice.E5 > 10 {
				myCountyData.GsNumE5 += 1
				myStateData.GsNumE5 += 1
				myStateData.AvgE5 += float64(myGasPrice.E5)
				myCountyData.AvgE5 += float64(myGasPrice.E5)
				/*
					if myCountyData.ID == 51 {
						gasStations := postCodeGasStationsMap[myPostCode]
						log.Printf("GsNumE5: %v, AvgE5: %v, E5: %v, Id: %v, GasStations: %v", myCountyData.GsNumE5, myCountyData.AvgE5, myGasPrice.E5, myCountyData.ID, len(gasStations))
					}
				*/
			}
			if myGasPrice.E10 > 10 {
				myCountyData.GsNumE10 += 1
				myStateData.GsNumE10 += 1
				myStateData.AvgE10 += float64(myGasPrice.E10)
				myCountyData.AvgE10 += float64(myGasPrice.E10)
			}
			if myGasPrice.Diesel > 10 {
				myCountyData.GsNumDiesel += 1
				myStateData.GsNumDiesel += 1
				myStateData.AvgDiesel += float64(myGasPrice.Diesel)
				myCountyData.AvgDiesel += float64(myGasPrice.Diesel)
			}
			idStateDataMap[int(myPostCodeLocation.StateData.ID)] = myStateData
			idCountyDataMap[int(myPostCodeLocation.CountyData.ID)] = myCountyData
		}
	}
	//log.Printf("e5Count: %v, e10Count: %v, dieselCount: %v", e5Count, e10Count, dieselCount)
	log.Printf("sums for postCodePostCodeLocationMap: %v", len(postCodeGasStationsMap))
	//divide by station count and save
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, myStateData := range idStateDataMap {
			if myStateData.GasStationNum > 0 {
				myStateData.AvgDiesel /= float64(myStateData.GsNumDiesel)
				myStateData.AvgE10 /= float64(myStateData.GsNumE10)
				myStateData.AvgE5 /= float64(myStateData.GsNumE5)
				tx.Save(&myStateData)
			}
		}
		for _, myCountyData := range idCountyDataMap {
			if myCountyData.GasStationNum > 0 {
				myCountyData.AvgDiesel /= float64(myCountyData.GsNumDiesel)
				myCountyData.AvgE10 /= float64(myCountyData.GsNumE10)
				myCountyData.AvgE5 /= float64(myCountyData.GsNumE5)
				tx.Save(&myCountyData)
			}
		}
		return nil
	})
	myDuration := time.Now().Sub(myStart)
	log.Printf("recalcCountyStatePrices finished for %v states and %v counties in %v.", len(idStateDataMap), len(idCountyDataMap), myDuration)
}

func createCodeTimeSliceBuckets(postCodePostCodeLocationMap map[int]pcmodel.PostCodeLocation, postCodeGasStationsMap map[string][]gsmodel.GasStation) map[string]map[time.Time][]gsmodel.GasPrice {
	gasStationIdGasPriceMap := createGasStationIdGasPriceArrayMap(&postCodeGasStationsMap, Day)
	postCodeTimeSliceBuckets := make(map[string]map[time.Time][]gsmodel.GasPrice)
	for _, myPostCodeLocation := range postCodePostCodeLocationMap {
		myPostCode := postcode.FormatPostCode(myPostCodeLocation.PostCode)
		for _, myGasStation := range postCodeGasStationsMap[myPostCode] {
			yesterday := time.Now().AddDate(0, 0, -1)
			//10 min slices
			timeSlices := make([]time.Time, 0)
			//10 min buckets
			timeSliceBuckets := make(map[time.Time][]gsmodel.GasPrice)
			for myTimeSlice := yesterday; myTimeSlice.Before(time.Now()); myTimeSlice = myTimeSlice.Add(time.Minute * 10) {
				timeSlices = append(timeSlices, myTimeSlice)
				timeSliceBuckets[myTimeSlice] = make([]gsmodel.GasPrice, 0)
			}
			for _, myStationPrice := range gasStationIdGasPriceMap[myGasStation.ID] {
				for i := 0; i < len(timeSlices)-1; i++ {
					if timeSlices[i].Before(myStationPrice.Date) && timeSlices[i+1].After(myStationPrice.Date) {
						timeSliceBuckets[timeSlices[i]] = append(timeSliceBuckets[timeSlices[i]], myStationPrice)
						break
					} else if i == len(timeSlices)-1 && timeSlices[i+1].Before(myStationPrice.Date) {
						timeSliceBuckets[timeSlices[i+1]] = append(timeSliceBuckets[timeSlices[i+1]], myStationPrice)
						break
					}
				}
			}
			for _, myTimeSlice := range timeSlices {
				if myBuckets, ok := postCodeTimeSliceBuckets[myPostCode]; ok {
					myBuckets[myTimeSlice] = append(myBuckets[myTimeSlice], timeSliceBuckets[myTimeSlice]...)
				} else {
					postCodeTimeSliceBuckets[myPostCode] = timeSliceBuckets
				}
			}
		}
	}
	return postCodeTimeSliceBuckets
}

func createTimeSliceBuckets(myCountyData pcmodel.CountyData, postCodePostCodeLocationMap map[int]pcmodel.PostCodeLocation, postCodeTimeSliceBuckets map[string]map[time.Time][]gsmodel.GasPrice) map[time.Time][]gsmodel.GasPrice {
	timeSliceBuckets := make(map[time.Time][]gsmodel.GasPrice)
	var myPostCodeLocations []pcmodel.PostCodeLocation
	for _, myPostCodeLocation := range postCodePostCodeLocationMap {
		if myPostCodeLocation.CountyDataID == myCountyData.ID && myCountyData.County != "Kreisfrei" {
			myPostCodeLocations = append(myPostCodeLocations, myPostCodeLocation)
		}
	}
	log.Printf("myPostCodeLocations: %v postcode: %v\n", len(myPostCodeLocations), myCountyData.County)
	for _, myPcLocation := range myPostCodeLocations {
		//log.Printf("postCodeLocations: %v\n", len(myPcLocation.CountyData.PostCodeLocations))
		myPostCode := postcode.FormatPostCode(myPcLocation.PostCode)
		if myBuckets, ok := postCodeTimeSliceBuckets[myPostCode]; ok {
			for mySlice, myNewPrices := range myBuckets {
				if myOldPrices, ok := timeSliceBuckets[mySlice]; ok {
					timeSliceBuckets[mySlice] = append(myOldPrices, myNewPrices...)
				} else {
					timeSliceBuckets[mySlice] = myNewPrices
				}
			}
		}
	}
	return timeSliceBuckets
}

func createTimeSlicesPriceSums(myCountyData pcmodel.CountyData, timeSliceBuckets map[time.Time][]gsmodel.GasPrice) map[time.Time]pcmodel.CountyTimeSlot {
	timeSlicesPriceSums := make(map[time.Time]pcmodel.CountyTimeSlot)
	for mySlice, myPrices := range timeSliceBuckets {
		var e5Sum, e10Sum, dieselSum, e5Num, e10Num, dieselNum, gsNum int64 = 0, 0, 0, 0, 0, 0, 0
		for _, myPrice := range myPrices {
			if myPrice.E5 > 10 {
				e5Num += 1
				e5Sum += int64(myPrice.E5)
			}
			if myPrice.E10 > 10 {
				e10Sum += int64(myPrice.E10)
				e10Num += 1
			}
			if myPrice.Diesel > 10 {
				dieselSum += int64(myPrice.Diesel)
				dieselNum += 1
			}
			if myPrice.E5 > 10 || myPrice.E10 > 10 || myPrice.Diesel > 10 {
				gsNum += 1
			}
		}
		if gsNum > 0 {
			myCountyTimeSlot := pcmodel.CountyTimeSlot{
				StartDate:    mySlice,
				CountyDataID: myCountyData.ID,
				CountyData:   myCountyData,
				GsAvgValues: pcmodel.GsAvgValues{
					GasStationNum: int(gsNum),
					GsNumDiesel:   int(dieselNum),
					GsNumE10:      int(e10Num),
					GsNumE5:       int(e5Num),
					AvgE5:         divZeroCheck(e5Sum, e5Num),
					AvgE10:        divZeroCheck(e10Sum, e10Num),
					AvgDiesel:     divZeroCheck(dieselSum, dieselNum),
				},
			}
			timeSlicesPriceSums[mySlice] = myCountyTimeSlot
		}
	}
	return timeSlicesPriceSums
}

func divZeroCheck(number int64, divisor int64) float64 {
	if divisor == 0 {
		return 0
	}
	return float64(number / divisor)
}

func CalcCountyTimeSlots() {
	log.Printf("calcCountyTimeSlots started.")
	myStart := time.Now()
	postCodePostCodeLocationMap, _, idCountyDataMap, postCodeGasStationsMap := createPostCodeGasStationMaps()
	postCodeTimeSliceBuckets := createCodeTimeSliceBuckets(postCodePostCodeLocationMap, postCodeGasStationsMap)
	//log.Printf("postCodeTimeSliceBuckets: %v\n", len(postCodeGasStationsMap))
	database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Unscoped().Where("1=1").Delete(&pcmodel.CountyTimeSlot{})
		return nil
	})
	//log.Printf("idCountyDataMap: %v\n", len(idCountyDataMap))
	//calc average changes in 10 min slots
	for _, myCountyData := range idCountyDataMap {
		//10 min buckets
		timeSliceBuckets := createTimeSliceBuckets(myCountyData, postCodePostCodeLocationMap, postCodeTimeSliceBuckets)
		//log.Printf("createTimeSliceBuckets: %v\n", len(timeSliceBuckets))
		timeSlicesPriceSums := createTimeSlicesPriceSums(myCountyData, timeSliceBuckets)
		database.DB.Transaction(func(tx *gorm.DB) error {
			for _, myCountyTimeSlots := range timeSlicesPriceSums {
				tx.Save(&myCountyTimeSlots)
			}
			return nil
		})
	}
	//store changes in countytimeslots
	myDuration := time.Now().Sub(myStart)
	log.Printf("calcCountyTimeSlots finished for %v counties in %v.", len(idCountyDataMap), myDuration)
}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).First(&myGasStation)
	return myGasStation
}

func FindPricesByStids(stids *[]string, resultLimit int, timeframe TimeFrame) []gsmodel.GasPrice {
	myGasPrice := findPricesByStids(stids, resultLimit, false, timeframe)
	return myGasPrice
}

func FindPricesByStidsDistinct(stids *[]string, resultLimit int, timeframe TimeFrame) []gsmodel.GasPrice {
	myGasPrice := findPricesByStids(stids, resultLimit, true, timeframe)
	return myGasPrice
}

func FindByPostCodes(postcodes []string) []gsmodel.GasStation {
	chunks := createChunks(&postcodes)
	var gasStations []gsmodel.GasStation
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, chunk := range chunks {
			var myGasStations []gsmodel.GasStation
			tx.Where("post_code IN ?", chunk).Find(&myGasStations)
			gasStations = append(gasStations, myGasStations...)
		}
		return nil
	})
	return gasStations
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
