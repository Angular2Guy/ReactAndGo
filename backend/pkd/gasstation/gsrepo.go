package gasstation

import (
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/gasstation/gsmodel"
	"fmt"

	"gorm.io/gorm"
)

func Start() {
	fmt.Println("Hello repo")
}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).Preload("GasPrices", func(db *gorm.DB) *gorm.DB {
		return db.Order("date DESC").Limit(20)
	}).First(&myGasStation)
	return myGasStation
}

func FindPricesByStid(stid string) []gsmodel.GasPrice {
	var myGasPrice []gsmodel.GasPrice
	database.DB.Where("stid = ?", stid).Order("date desc").First(&myGasPrice)
	return myGasPrice
}
