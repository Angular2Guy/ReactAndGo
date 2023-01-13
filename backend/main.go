package main

import (
	"angular-and-go/pkd/config"
	"angular-and-go/pkd/contr"
	"angular-and-go/pkd/cron"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/messaging"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateDB()
	messaging.Start()
	cron.Start()
}

func main() {
	contr.Start()
}
