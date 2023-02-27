package notification

import (
	"log"
	"react-and-go/pkd/appuser"
	"react-and-go/pkd/gasstation/gsmodel"
)

type gasStationWithPrice struct {
	gasStation gsmodel.GasStation
	gasPrice   gsmodel.GasPrice
}

func SendNotifications(gasStationIDToGasPriceMap map[string]gsmodel.GasPrice, gasStations []gsmodel.GasStation) {
	gasStationWithPricesMap := make(map[string]gasStationWithPrice)
	for _, gasStation := range gasStations {
		myGasStationWithPrice := gasStationWithPrice{}
		myGasStationWithPrice.gasPrice = gasStationIDToGasPriceMap[gasStation.ID]
		myGasStationWithPrice.gasStation = gasStation
		gasStationWithPricesMap[gasStation.ID] = myGasStationWithPrice
	}
	allAppUsers := appuser.FindAllUsers()
	for _, appUser := range allAppUsers {
		for _, myGasStationWithPrice := range gasStationWithPricesMap {
			//Target price reached?
			if appUser.TargetDiesel >= myGasStationWithPrice.gasPrice.Diesel || appUser.TargetE10 >= myGasStationWithPrice.gasPrice.E10 || appUser.TargetE5 >= myGasStationWithPrice.gasPrice.E5 {
				//Distance match?
				distance, _ := myGasStationWithPrice.gasStation.CalcDistanceBearing(appUser.Latitude, appUser.Longitude)
				if appUser.SearchRadius >= distance {
					log.Printf("Match found: %v %v %v %v %v %v\n", distance, myGasStationWithPrice.gasStation.Brand, myGasStationWithPrice.gasStation.Place,
						myGasStationWithPrice.gasPrice.Diesel, myGasStationWithPrice.gasPrice.E10, myGasStationWithPrice.gasPrice.E5)
				}
			}
		}
	}
}
