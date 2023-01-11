package cron

import (
	gsclient "angular-and-go/pkd/contr/client"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

type CircleCenter struct {
	Latitude  float64
	Longitude float64
}

var HamburgAndSH = [...]CircleCenter{
	{
		Latitude:  54.824158,
		Longitude: 8.346131,
	},
	{
		Latitude:  54.715297,
		Longitude: 8.775641,
	},
	{
		Latitude:  54.661861,
		Longitude: 9.180214,
	},
	{
		Latitude:  54.677340,
		Longitude: 9.743868,
	},
	{
		Latitude:  54.298884,
		Longitude: 8.743990,
	},
	{
		Latitude:  54.308298,
		Longitude: 9.317139,
	},
	{
		Latitude:  54.306721,
		Longitude: 9.792173,
	},
	{
		Latitude:  54.280894,
		Longitude: 10.247840,
	},
	{
		Latitude:  54.333907,
		Longitude: 10.987011,
	},
	{
		Latitude:  54.019711,
		Longitude: 10.643870,
	},
	{
		Latitude:  53.889138,
		Longitude: 10.020025,
	},
	{
		Latitude:  53.913517,
		Longitude: 9.572239,
	},
	{
		Latitude:  53.928135,
		Longitude: 9.042212,
	},
	{
		Latitude:  53.648308,
		Longitude: 10.580193,
	},
	{
		Latitude:  53.473590,
		Longitude: 10.277897,
	},
	{
		Latitude:  53.522599,
		Longitude: 9.800100,
	},
}

var requestCounter int64 = 0

func Start() {
	var apikeys [3]string
	for index, _ := range apikeys {
		apikeys[index] = os.Getenv(fmt.Sprintf("APIKEY%v", index))
	}
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().At("01:07").Do(func() {
		gsclient.UpdateGasStations(nil)
	})
	scheduler.Every(2).Minutes().Tag("prices").Do(updatePriceRegion, HamburgAndSH)
	/*
		for _, value := range HamburgAndSH {
			fmt.Printf("Lag: %f Lng: %f Rad: %f\n", value.Latitude, value.Longitude, 25.0)
			scheduler.Every(2).Minutes().Tag("prices").Do(gsclient.UpdateGsPrices, value.Latitude, value.Longitude, 25.0)
		}
	*/
	scheduler.StartAsync()
}

func updatePriceRegion(regionCircleCenters [16]CircleCenter, apikeys [3]string) {
	apikeyIndex := 0
	for _, value := range HamburgAndSH {
		time.Sleep(6 * time.Second)
		gsclient.UpdateGsPrices(value.Latitude, value.Longitude, 25.0, apikeys[apikeyIndex])
		requestCounter += 1
		if requestCounter%45 == 0 {
			if apikeyIndex < 2 {
				apikeyIndex += 1
			} else {
				apikeyIndex = 0
			}
		}
		log.Printf("Request %v, ApikeyIndex: %v", requestCounter, apikeyIndex)
	}
}
