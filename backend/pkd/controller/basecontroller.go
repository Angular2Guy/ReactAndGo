package controller

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	gsclient "react-and-go/pkd/controller/client"
	token "react-and-go/pkd/token"
	"strings"

	"github.com/gin-gonic/gin"
)

func Start(embeddedFiles fs.FS) {
	router := gin.Default()
	router.POST("/appuser/signin", postSignin)
	router.POST("/appuser/login", postLogin)
	router.GET("/appuser/logout", token.CheckToken, getLogout)
	router.GET("/appuser/location", token.CheckToken, getLocation)
	router.GET("/appuser/refreshtoken", token.CheckToken, getRefreshToken)
	router.POST("/appuser/locationradius", token.CheckToken, postUserLocationRadius)
	router.POST("/appuser/targetprices", token.CheckToken, postTargetPrices)
	router.GET("/config/updategs", token.CheckToken, gsclient.UpdateGasStations)
	router.GET("/config/updatepc", token.CheckToken, getPostCodeCoordinates)
	router.GET("/gasprice/:id", token.CheckToken, getGasPriceByGasStationId)
	router.GET("/gasstation/:id", token.CheckToken, getGasStationById)
	router.POST("/gasstation/search/place", token.CheckToken, searchGasStationPlace)
	router.POST("/gasstation/search/location", token.CheckToken, searchGasStationLocation)
	router.GET("/usernotification/new/:useruuid", token.CheckToken, getNewUserNotifications)
	router.GET("/usernotification/current/:useruuid", token.CheckToken, getCurrentUserNotifications)
	router.StaticFS("/public", http.FS(embeddedFiles))
	//router.Static("/public", "./public")
	router.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/public") })
	absolutePathKeyFile := strings.TrimSpace(os.Getenv("ABSOLUTE_PATH_KEY_FILE"))
	absolutePathCertFile := strings.TrimSpace(os.Getenv("ABSOLUTE_PATH_CERT_FILE"))
	if len(absolutePathCertFile) < 2 || len(absolutePathKeyFile) < 2 {
		router.Run() // listen and serve on 0.0.0.0:3000
	} else {
		log.Fatal(router.RunTLS("0.0.0.0", absolutePathCertFile, absolutePathKeyFile))
	}
}
