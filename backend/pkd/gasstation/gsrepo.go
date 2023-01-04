package gasstation

import (
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/gasstation/gsmodel"
	"fmt"
)

func Start() {
	fmt.Println("Hello repo")
}

func FindById(id string) gsmodel.GasStation {
	var myGasStation gsmodel.GasStation
	database.DB.Where("id = ?", id).First(&myGasStation)
	//.Preload("GasPrices")
	//.Joins("left join (select * from gas_station_information_history order by date desc) as gsih on gsih.stid")
	return myGasStation
}

func FindPricesByStid(stid string) []gsmodel.GasPrice {
	var myGasPrice []gsmodel.GasPrice
	database.DB.Where("stid = ?", stid).Order("date desc").First(&myGasPrice)
	return myGasPrice
}
