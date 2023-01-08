package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func Start() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Minutes().Do(func() {
		log.Default().Printf("Import Prices:")
	})
	scheduler.StartAsync()
}
