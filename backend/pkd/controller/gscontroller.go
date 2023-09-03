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
package controller

import (
	"log"
	"net/http"
	gsbody "react-and-go/pkd/controller/gsmodel"
	"react-and-go/pkd/gasstation"

	"github.com/gin-gonic/gin"
)

func getGasPriceByGasStationId(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	gsEntity := gasstation.FindPricesByStid(gasstationId)
	c.JSON(http.StatusOK, gsEntity)
}

func getGasStationById(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	gsEntity := gasstation.FindById(gasstationId)
	c.JSON(http.StatusOK, gsEntity)
}

func searchGasStationPlace(c *gin.Context) {
	var searchPlaceBody gsbody.SearchPlaceBody
	if err := c.Bind(&searchPlaceBody); err != nil {
		log.Printf("searchGasStationPlace: %v", err.Error())
	}
	gsEntity := gasstation.FindBySearchPlace(searchPlaceBody)
	c.JSON(http.StatusOK, gsEntity)
}

func searchGasStationLocation(c *gin.Context) {
	//jsonData, err := ioutil.ReadAll(c.Request.Body)
	//fmt.Printf("Json: %v, Err: %v", string(jsonData), err)
	var searchLocationBody gsbody.SearchLocation
	if err := c.Bind(&searchLocationBody); err != nil {
		log.Printf("searchGasStationLocation: %v", err.Error())
	}
	//fmt.Printf("Lat: %v, Lng: %v\n", searchLocationBody.Latitude, searchLocationBody.Longitude)
	gsEntity := gasstation.FindBySearchLocation(searchLocationBody)
	c.JSON(http.StatusOK, gsEntity)
}

func getRecalcAvgs(c *gin.Context) {
	gasstation.ReCalcCountyStatePrices()
	c.JSON(http.StatusOK, "Done.")
}
