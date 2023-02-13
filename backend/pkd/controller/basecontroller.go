package controller

import (
	gsclient "angular-and-go/pkd/controller/client"
	token "angular-and-go/pkd/token"

	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	router.POST("/appuser/signin", postSignin)
	router.POST("/appuser/login", postLogin)
	router.GET("/config/updategs", token.CheckToken, gsclient.UpdateGasStations)
	router.GET("/config/updatepc", token.CheckToken, getPostCodeCoordinates)
	router.GET("/gasprice/:id", token.CheckToken, getGasPriceByGasStationId)
	router.GET("/gasstation/:id", token.CheckToken, getGasStationById)
	router.POST("/gasstation/search/place", token.CheckToken, searchGasStationPlace)
	router.POST("/gasstation/search/location", token.CheckToken, searchGasStationLocation)
	router.PUT("/posts/:id", postsUpdate)
	router.DELETE("/posts/:id", postsDelete)
	router.Static("/static", "./static")
	router.NoRoute(func(c *gin.Context) { c.Redirect(302, "/static") })
	router.Run() // listen and serve on 0.0.0.0:3000
}
