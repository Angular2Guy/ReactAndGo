package main

import (
	"angular-and-go/pkd/config"
	"angular-and-go/pkd/contr"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/gasstation"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateDB()
}

func main() {
	gasstation.Start()
	contr.Start()
}
