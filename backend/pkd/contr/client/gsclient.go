package gsclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateGsPrices(c *gin.Context) {
	//func UpdateGsPrices(latitude float64, longitude float64, radiusKM float64) {
	var latitude = 52.521
	var longitude = 13.438
	var radiusKM = 10.0
	var queryUrl = fmt.Sprintf("https://creativecommons.tankerkoenig.de/json/list.php?lat=%f&lng=%f&rad=%f&sort=dist&type=all&apikey=00000000-0000-0000-0000-000000000002", latitude, longitude, radiusKM)
	response, err := http.Get(queryUrl)
	if err != nil {
		log.Fatalf("Request failed: %v", err.Error())
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Read body failed: %v", body)
	}
	defer response.Body.Close()
	fmt.Printf("Body received: %v", string(body))
}
