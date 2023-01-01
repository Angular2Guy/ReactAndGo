package main

import (
	"angular-and-go/pkd/config"
	contr "angular-and-go/pkd/contr"
	repo "angular-and-go/pkd/gas-stations"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	repo.Start()
	contr.Start()
}
