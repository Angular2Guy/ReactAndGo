package controller

import (
	gsbody "angular-and-go/pkd/controller/gsmodel"
	"angular-and-go/pkd/gasstation"
	"log"
	"net/http"

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

func postsUpdate(c *gin.Context) {

}

func postsDelete(c *gin.Context) {

}
