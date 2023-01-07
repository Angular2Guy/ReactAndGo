package gsclient

import (
	"angular-and-go/pkd/gasstation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func UpdateGsPrices(c *gin.Context) {
	//func UpdateGsPrices(latitude float64, longitude float64, radiusKM float64) {
	var latitude = 52.521
	var longitude = 13.438
	var radiusKM = 10.0
	var queryUrl = fmt.Sprintf("https://creativecommons.tankerkoenig.de/json/list.php?lat=%f&lng=%f&rad=%f&sort=dist&type=all&apikey=00000000-0000-0000-0000-000000000002", latitude, longitude, radiusKM)
	response, err := http.Get(queryUrl)
	if err != nil {
		log.Fatalf("Request failed: %v\n", err.Error())
	}
	defer response.Body.Close()
	var myGsResponse gsResponse
	if err := json.NewDecoder(response.Body).Decode(&myGsResponse); err != nil {
		log.Fatalf("Json decode failed: %v", err.Error())
	}
	stationPricesMap := make(map[string]gasstation.GasStationPrices)
	for _, value := range myGsResponse.Stations {
		stationPricesMap[value.Id] = gasstation.GasStationPrices{GasStationID: value.Id, E5: int(value.E5 * 1000), E10: int(value.E10 * 1000), Diesel: int(value.Diesel * 1000), Date: time.Now()}
	}
	var gasPriceUpdates []gasstation.GasStationPrices
	for _, value := range stationPricesMap {
		gasPriceUpdates = append(gasPriceUpdates, value)
	}
	gasstation.UpdatePrice(gasPriceUpdates)
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
	result, _ := json.Marshal(myGsResponse)
	fmt.Printf("Json: %v\n", string(result))
}
