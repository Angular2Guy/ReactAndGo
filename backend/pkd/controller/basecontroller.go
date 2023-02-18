package controller

import (
	gsclient "angular-and-go/pkd/controller/client"
	token "angular-and-go/pkd/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.POST("/appuser/signin", postSignin)
	router.POST("/appuser/login", postLogin)
	router.GET("/appuser/location", token.CheckToken, getLocation)
	router.PUT("/appuser/locationradius", putUserLocationRadius)
	router.PUT("/appuser/targetprices", putTargetPrices)
	router.GET("/config/updategs", token.CheckToken, gsclient.UpdateGasStations)
	router.GET("/config/updatepc", token.CheckToken, getPostCodeCoordinates)
	router.GET("/gasprice/:id", token.CheckToken, getGasPriceByGasStationId)
	router.GET("/gasstation/:id", token.CheckToken, getGasStationById)
	router.POST("/gasstation/search/place", token.CheckToken, searchGasStationPlace)
	router.POST("/gasstation/search/location", token.CheckToken, searchGasStationLocation)
	router.PUT("/posts/:id", postsUpdate)
	router.Static("/static", "./static")
	router.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/static") })
	router.Run() // listen and serve on 0.0.0.0:3000
}
