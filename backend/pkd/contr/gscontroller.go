package contr

import (
	gsclient "angular-and-go/pkd/contr/client"
	gsbody "angular-and-go/pkd/contr/model"
	"angular-and-go/pkd/gasstation"

	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.POST("/posts", postsCreate)
	router.GET("/clienttest", gsclient.UpdateGsPrices)
	router.GET("/gasprice/:id", getGasPriceByGasStationId)
	router.GET("/gasstation/:id", getGasStationById)
	router.POST("/gasstation/search/place", searchGasStationPlace)
	router.POST("/gasstation/search/location", searchGasStationLocation)
	router.PUT("/posts/:id", postsUpdate)
	router.DELETE("/posts/:id", postsDelete)
	router.Static("/static", "./static")
	router.NoRoute(func(c *gin.Context) { c.Redirect(302, "/static") })
	router.Run() // listen and serve on 0.0.0.0:3000
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
