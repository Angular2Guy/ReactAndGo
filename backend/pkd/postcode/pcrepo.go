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
package postcode

import (
	"fmt"
	"log"
	"react-and-go/pkd/database"
	pcmodel "react-and-go/pkd/postcode/pcmodel"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PostCodeData struct {
	Label           string
	PostCode        int32
	Population      int32
	SquareKM        float32
	CenterLongitude float64
	CenterLatitude  float64
}

func FindLocation(locationStr string) []pcmodel.PostCodeLocation {
	result := []pcmodel.PostCodeLocation{}
	database.DB.Where("lower(label) like ?", fmt.Sprintf("%%%v%%", strings.ToLower(strings.TrimSpace(locationStr)))).Limit(20).Find(&result)
	//log.Printf("Select: %v failed. %v", fmt.Sprintf("%%%v%%", strings.ToLower(strings.TrimSpace(locationStr))), err)
	return result
}

func ImportPostCodeData(postCodeData []PostCodeData) {
	postCodeLocations := mapToPostCodeLocation(postCodeData)
	var oriPostCodeLocations []pcmodel.PostCodeLocation
	database.DB.Find(&oriPostCodeLocations)
	var myCountyData pcmodel.CountyData
	var myStateData pcmodel.StateData
	postCodeLocationsMap := make(map[int32]pcmodel.PostCodeLocation)
	for _, oriPostCodeLocation := range oriPostCodeLocations {
		postCodeLocationsMap[oriPostCodeLocation.PostCode] = oriPostCodeLocation
	}
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, postCodeLocation := range postCodeLocations {
			oriPostCodeLocation, exists := postCodeLocationsMap[postCodeLocation.PostCode]
			if exists {
				oriPostCodeLocation.Label = postCodeLocation.Label
				oriPostCodeLocation.PostCode = postCodeLocation.PostCode
				oriPostCodeLocation.Population = postCodeLocation.Population
				oriPostCodeLocation.SquareKM = postCodeLocation.SquareKM
				oriPostCodeLocation.CenterLongitude = postCodeLocation.CenterLongitude
				oriPostCodeLocation.CenterLatitude = postCodeLocation.CenterLatitude
				tx.Save(&oriPostCodeLocation)
			} else {
				tx.Save(&myCountyData)
				tx.Save(&myStateData)
				postCodeLocation.CountyData = myCountyData
				postCodeLocation.CountyDataID = myCountyData.ID
				postCodeLocation.StateData = myStateData
				postCodeLocation.StateDataID = myStateData.ID
				tx.Save(&postCodeLocation)
			}
		}
		return nil
	})
	log.Printf("PostCodeLocations saved: %v\n", len(postCodeLocations))
}

func UpdateStatesCounties(plzToState map[string]string, plzToCounty map[string]string) {
	var pcLocations []pcmodel.PostCodeLocation
	database.DB.Preload("StateData").Preload("CountyData").Find(&pcLocations)
	stateMap := make(map[string]*pcmodel.StateData)
	countyMap := make(map[string]*pcmodel.CountyData)
	//log.Printf("%d pcLocations.", len(pcLocations))
	//log.Printf("%s, %s", plzToCounty[FormatPostCode(1159)], plzToState[FormatPostCode(1159)])
	database.DB.Transaction(func(tx *gorm.DB) error {
		for _, pcLocation := range pcLocations {
			if &pcLocation.CountyData.County == nil || pcLocation.CountyData.County == "" {
				myCountyData := pcmodel.CountyData{}
				if mapValue, ok := countyMap[plzToCounty[FormatPostCode(pcLocation.PostCode)]]; ok {
					myCountyData = *mapValue
				} else {
					countyMap[plzToCounty[FormatPostCode(pcLocation.PostCode)]] = &myCountyData
					myCountyData.County = plzToCounty[FormatPostCode(pcLocation.PostCode)]
				}
				tx.Save(&myCountyData)
				pcLocation.CountyData = myCountyData
				pcLocation.CountyDataID = myCountyData.ID
			}
			if &pcLocation.StateData.State == nil || pcLocation.StateData.State == "" {
				myStateData := pcmodel.StateData{}
				if myMapValue, ok := stateMap[plzToState[FormatPostCode(pcLocation.PostCode)]]; ok {
					myStateData = *myMapValue
				} else {
					stateMap[plzToState[FormatPostCode(pcLocation.PostCode)]] = &myStateData
					myStateData.State = plzToState[FormatPostCode(pcLocation.PostCode)]
				}
				tx.Save(&myStateData)
				pcLocation.StateData = myStateData
				pcLocation.StateDataID = myStateData.ID
			}
			pcLocation.CountyData.County = plzToCounty[FormatPostCode(pcLocation.PostCode)]
			pcLocation.StateData.State = plzToState[FormatPostCode(pcLocation.PostCode)]
			tx.Save(&pcLocation)
		}
		myCountyData := pcmodel.CountyData{}
		tx.Where("county = ? or county is null", "").Delete(&myCountyData)
		myStateData := pcmodel.StateData{}
		tx.Where("state = ? or state is null", "").Delete(&myStateData)
		return nil
	})
	log.Printf("UpdateStatesCounties updated: %v\n", len(pcLocations))
}

func FindByPlzs(plzs []string) *[]pcmodel.PostCodeLocation {
	var pcLocations []pcmodel.PostCodeLocation
	var plzInts = plzsToPlzInts(plzs)
	database.DB.Where("post_code in ?", plzInts).Preload("StateData").Preload("CountyData").Find(&pcLocations)
	return &pcLocations
}

func plzsToPlzInts(plzs []string) []int {
	var plzInts []int
	for _, myPlz := range plzs {
		myPlzInt, err := strconv.Atoi(myPlz)
		if err != nil {
			log.Printf("Failed to parse: %v", myPlz)
		} else {
			plzInts = append(plzInts, myPlzInt)
		}
	}
	return plzInts
}

func FormatPostCode(postCode int32) string {
	pcStr := strconv.Itoa(int(postCode))
	for len(pcStr) < 5 {
		pcStr = "0" + pcStr
	}
	return pcStr
}

func mapToPostCodeLocation(postCodeData []PostCodeData) []pcmodel.PostCodeLocation {
	result := []pcmodel.PostCodeLocation{}
	for _, myPostCodeData := range postCodeData {
		myPostCodeLocation := pcmodel.PostCodeLocation{}
		myPostCodeLocation.Label = myPostCodeData.Label
		myPostCodeLocation.PostCode = myPostCodeData.PostCode
		myPostCodeLocation.Population = myPostCodeData.Population
		myPostCodeLocation.SquareKM = myPostCodeData.SquareKM
		myPostCodeLocation.CenterLongitude = myPostCodeData.CenterLongitude
		myPostCodeLocation.CenterLatitude = myPostCodeData.CenterLatitude
		result = append(result, myPostCodeLocation)
	}
	return result
}
