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
package gsclient

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"react-and-go/pkd/gasstation"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type gsResponse struct {
	Ok       bool         `json:"ok"`
	License  string       `json:"license"`
	Data     string       `json:"data"`
	Status   string       `json:"status"`
	Stations []gsStations `json:"stations"`
}

type gsStations struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Street      string  `json:"street"`
	Place       string  `json:"place"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Dist        float64 `json:"dist"`
	Diesel      float64 `json:"diesel"`
	E5          float64 `json:"e5"`
	E10         float64 `json:"e10"`
	IsOpen      bool    `json:"isOpen"`
	HouseNumber string  `json:"houseNumber"`
	PostCode    int     `json:"postCode"`
}

func UpdateGasStations(c *gin.Context) {
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	//url := fmt.Sprintf("https://dev.azure.com/tankerkoenig/362e70d1-bafa-4cf7-a346-1f3613304973/_apis/git/repositories/0d6e7286-91e4-402c-af56-fa75be1f223d/items?path=/stations/%04d/%02d/%04d-%02d-%02d-stations.csv&versionDescriptor%5BversionOptions%5D=0&versionDescriptor%5BversionType%5D=0&versionDescriptor%5Bversion%5D=master&resolveLfs=true&%24format=octetStream&api-version=5.0&download=true", year, month, year, month, day)
	url := fmt.Sprintf("https://dev.azure.com/tankerkoenig/362e70d1-bafa-4cf7-a346-1f3613304973/_apis/git/repositories/0d6e7286-91e4-402c-af56-fa75be1f223d/Items?path=/stations/%04d/%02d/%04d-%02d-%02d-stations.csv"+
		"&recursionLevel=0&includeContentMetadata=true&versionDescriptor.version=master&versionDescriptor.versionOptions=0&versionDescriptor.versionType=0&includeContent=true&resolveLfs=true", year, month, year, month, day)
	//fmt.Printf("Url: %v\n", url)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Request failed: %v\n", url)
	}
	defer response.Body.Close()
	reader := csv.NewReader(response.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Println("Cannot read request body:", err)
	}
	gasStationImports := convertCsvToGasStationImports(&rows)
	log.Printf("Result: %v\n", len(gasStationImports))
	//log.Default().Fatalf("Result: GasStationImport {\nUuid: %v\nStationName: %v\nBrand: %v\nStreet: %v\nHouseNumber: %v\nPostCode: %v\nCity: %v\nLatitude: %v\nLongitude: %v\nFirstActive: %v\nOpeningTimesJson: %v\n",
	//	gasStationImports[0].Uuid, gasStationImports[0].StationName, gasStationImports[0].Brand, gasStationImports[0].Street, gasStationImports[0].HouseNumber, gasStationImports[0].PostCode, gasStationImports[0].City,
	//	gasStationImports[0].Latitude, gasStationImports[0].Longitude, gasStationImports[0].FirstActive, gasStationImports[0].OpeningTimesJson)
	gasstation.UpdateGasStations(&gasStationImports)
}

func convertCsvToGasStationImports(rows *[][]string) []gasstation.GasStationImport {
	var result []gasstation.GasStationImport
	for _, row := range *rows {
		//ignore header
		if strings.ToLower(row[0]) == "uuid" {
			continue
		}
		lat, _ := strconv.ParseFloat(row[7], 64)
		lng, _ := strconv.ParseFloat(row[8], 64)
		//log.Default().Printf("Date: %v", strings.TrimSpace(row[9]))
		firstActive, _ := time.Parse("2006-02-01 15:04:05-07", strings.TrimSpace(row[9]))
		//log.Default().Printf("GoDate: %v", firstActive.UTC())
		gsImport := gasstation.GasStationImport{Uuid: row[0],
			StationName:      row[1],
			Brand:            row[2],
			Street:           row[3],
			HouseNumber:      row[4],
			PostCode:         row[5],
			City:             row[6],
			Latitude:         lat,
			Longitude:        lng,
			FirstActive:      firstActive,
			OpeningTimesJson: row[10],
		}
		result = append(result, gsImport)
	}
	return result
}

func UpdateGsPrices1(c *gin.Context) {
	var latitude = 52.521
	var longitude = 13.438
	var radiusKM = 10.0
	apikey := os.Getenv("APIKEY1")
	UpdateGsPrices(latitude, longitude, radiusKM, apikey)
}

func UpdateGsPrices(latitude float64, longitude float64, radiusKM float64, apikey string) error {
	fmt.Printf("Price requested Latitude: %f Longitude: %f radiusKM:: %f\n", latitude, longitude, radiusKM)
	var queryUrl = fmt.Sprintf("https://creativecommons.tankerkoenig.de/json/list.php?lat=%f&lng=%f&rad=%f&sort=dist&type=all&apikey=%v", latitude, longitude, radiusKM, strings.TrimSpace(apikey))
	//log.Printf("Url: %v", queryUrl)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(queryUrl)
	if err != nil {
		log.Printf("Request failed: %v\n", err.Error())
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 {
		log.Printf("Response status: %v\n", response.Status)
		err := fmt.Errorf("response status: %v", response.Status)
		return err
	}
	var myGsResponse gsResponse
	if err := json.NewDecoder(response.Body).Decode(&myGsResponse); err != nil {
		log.Printf("Json decode failed: %v\n", err.Error())
		return err
	}
	stationPricesMap := make(map[string]gasstation.GasStationPrices)
	for _, value := range myGsResponse.Stations {
		stationPricesMap[value.Id] = gasstation.GasStationPrices{GasStationID: value.Id, E5: int(value.E5 * 1000), E10: int(value.E10 * 1000), Diesel: int(value.Diesel * 1000), Timestamp: time.Now()}
	}
	var gasPriceUpdates []gasstation.GasStationPrices
	for _, value := range stationPricesMap {
		gasPriceUpdates = append(gasPriceUpdates, value)
	}
	fmt.Printf("Number of Price updates: %v\n", len(gasPriceUpdates))
	gasstation.UpdatePrice(&gasPriceUpdates)
	/*
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Read body failed: %v\n", body)
		}
		//fmt.Printf("Body received: %v", string(body))
		err1 := json.Unmarshal(body, &myGsResponse)
		if err1 != nil {
			log.Fatalf("Error: %v\n", err1.Error())
		}
	*/
	/*
		result, _ := json.Marshal(myGsResponse)
		fmt.Printf("Json: %v\n", string(result))
	*/
	return nil
}
