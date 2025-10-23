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
	"io/fs"
	"log"
	"net/http"
	"os"
	gsclient "react-and-go/pkd/controller/client"
	token "react-and-go/pkd/token"
	"strconv"
	"strings"

	"github.com/angular2guy/go-actuator"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Start(embeddedFiles fs.FS) {
	apiBase := "/api"
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.POST(apiBase+"/appuser/signin", postSignin)
	router.POST(apiBase+"/appuser/login", postLogin)
	router.GET(apiBase+"/appuser/logout", token.CheckToken, getLogout)
	router.GET(apiBase+"/appuser/location", token.CheckToken, getLocation)
	router.GET(apiBase+"/appuser/refreshtoken", token.CheckToken, getRefreshToken)
	router.POST(apiBase+"/appuser/locationradius", token.CheckToken, postUserLocationRadius)
	router.POST(apiBase+"/appuser/targetprices", token.CheckToken, postTargetPrices)
	router.GET(apiBase+"/config/updategs", token.CheckToken, gsclient.UpdateGasStations)
	router.GET(apiBase+"/config/updatepc", token.CheckToken, getPostCodeCoordinates)
	router.GET(apiBase+"/config/updatestatescounties", token.CheckToken, getStateCountyData)
	router.GET(apiBase+"/config/recalcAvgs", token.CheckToken, getRecalcAvgs)
	router.GET(apiBase+"/gasprice/:id", token.CheckToken, getGasPriceByGasStationId)
	router.GET(apiBase+"/gasstation/:id", token.CheckToken, getGasStationById)
	router.GET(apiBase+"/gasprice/avgs/:postcode", token.CheckToken, getAveragePrices)
	router.POST(apiBase+"/gasstation/search/place", token.CheckToken, searchGasStationPlace)
	router.POST(apiBase+"/gasstation/search/location", token.CheckToken, searchGasStationLocation)
	router.GET(apiBase+"/usernotification/new/:useruuid", token.CheckToken, getNewUserNotifications)
	router.GET(apiBase+"/usernotification/current/:useruuid", token.CheckToken, getCurrentUserNotifications)
	router.GET(apiBase+"/postcode/countytimeslots/:postcode", token.CheckToken, getCountyDataByIdWithTimeSlots)
	router.GET(apiBase+"/gasstation/countytimeslots/recalc", token.CheckToken, getRecalcTimeSlots)

	myPort := strings.TrimSpace(os.Getenv("PORT"))
	portNum, err := strconv.ParseInt(myPort, 10, 0)
	if err != nil {
		log.Fatal("Failed to parse port to int: " + myPort)
	}
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Port: int(portNum)})
	ginActuatorHandler := func(ctx *gin.Context) {
		actuatorHandler(ctx.Writer, ctx.Request)
	}
	router.GET("/actuator/*endpoint", ginActuatorHandler)

	router.StaticFS("/public", http.FS(embeddedFiles))
	//router.Static("/public", "./public")
	router.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/public") })
	absolutePathKeyFile := strings.TrimSpace(os.Getenv("ABSOLUTE_PATH_KEY_FILE"))
	absolutePathCertFile := strings.TrimSpace(os.Getenv("ABSOLUTE_PATH_CERT_FILE"))
	if len(absolutePathCertFile) < 2 || len(absolutePathKeyFile) < 2 || len(myPort) < 2 {
		router.Run() // listen and serve on 0.0.0.0:3000
	} else {
		log.Fatal(router.RunTLS(":"+myPort, absolutePathCertFile, absolutePathKeyFile))
	}
}
