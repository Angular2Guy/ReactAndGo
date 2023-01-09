package cron

import (
	gsclient "angular-and-go/pkd/contr/client"
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().At("01:07").Do(func() {
		gsclient.UpdateGasStations(nil)
	})
	scheduler.StartAsync()
}
