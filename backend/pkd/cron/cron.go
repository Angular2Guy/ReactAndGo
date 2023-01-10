package cron

import (
	gsclient "angular-and-go/pkd/contr/client"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type CircleCenter struct {
	Latitude  float64
	Longitude float64
}

var HamburgAndSH = []CircleCenter{
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

func Start() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().At("01:07").Do(func() {
		gsclient.UpdateGasStations(nil)
	})
	for index, value := range HamburgAndSH {
		scheduler.Every(2).Minutes().Tag(fmt.Sprintf("tag%v", index)).Do(func() {
			gsclient.UpdateGsPrices(value.Latitude, value.Longitude, 25.0)
		})
	}
	scheduler.StartAsync()
	for index, _ := range HamburgAndSH {
		scheduler.RunByTagWithDelay(fmt.Sprintf("tag%v", index), time.Duration(index*6)*time.Second)
	}
}
