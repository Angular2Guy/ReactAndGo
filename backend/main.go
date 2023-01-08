package main

import (
	"angular-and-go/pkd/config"
	"angular-and-go/pkd/contr"
	"angular-and-go/pkd/cron"
	"angular-and-go/pkd/database"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateDB()
	cron.Start()
}

func main() {
	contr.Start()
}
