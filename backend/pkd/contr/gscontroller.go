package contr

import (
	gsbody "angular-and-go/pkd/contr/model"
	"angular-and-go/pkd/gasstation"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.POST("/posts", postsCreate)
	r.GET("/gasprice/:id", getGasPriceByGasStationId)
	r.GET("/gasstation/:id", getGasStationById)
	r.POST("/gasstation/search/place", searchGasStationPlace)
	r.POST("/gasstation/search/location", searchGasStationLocation)
	r.PUT("/posts/:id", postsUpdate)
	r.DELETE("posts/:id", postsDelete)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func postsCreate(c *gin.Context) {

}

func getGasPriceByGasStationId(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	gsEntity := gasstation.FindPricesByStid(gasstationId)
	c.JSON(200, gsEntity)
}

func getGasStationById(c *gin.Context) {
	gasstationId := c.Params.ByName("id")
	gsEntity := gasstation.FindById(gasstationId)
	c.JSON(200, gsEntity)
}

func searchGasStationPlace(c *gin.Context) {
	var searchPlaceBody gsbody.SearchPlaceBody
	c.Bind(&searchPlaceBody)
	gsEntity := gasstation.FindBySearchPlace(searchPlaceBody)
	c.JSON(200, gsEntity)
}

func searchGasStationLocation(c *gin.Context) {
	//jsonData, err := ioutil.ReadAll(c.Request.Body)
	//fmt.Printf("Json: %v, Err: %v", string(jsonData), err)
	var searchLocationBody gsbody.SearchLocation
	c.Bind(&searchLocationBody)
	//fmt.Printf("Lat: %v, Lng: %v\n", searchLocationBody.Latitude, searchLocationBody.Longitude)
	gsEntity := gasstation.FindBySearchLocation(searchLocationBody)
	c.JSON(200, gsEntity)
}

func postsUpdate(c *gin.Context) {

}

func postsDelete(c *gin.Context) {

}
