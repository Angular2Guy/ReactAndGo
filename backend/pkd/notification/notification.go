package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"react-and-go/pkd/appuser"
	"react-and-go/pkd/gasstation/gsmodel"
	"time"
)

type gasStationWithPrice struct {
	gasStation gsmodel.GasStation
	gasPrice   gsmodel.GasPrice
}

type NotificationData struct {
	GasStationID string
	StationName  string
	Brand        string
	Street       string
	Place        string
	HouseNumber  string
	PostCode     string
	Latitude     float64
	Longitude    float64
	E5           int
	E10          int
	Diesel       int
	Timestamp    time.Time
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
	myNotificationMsgs := []NotificationMsg{}
	for _, appUser := range allAppUsers {
		gsMatches := []gasStationWithPrice{}
		for _, myGasStationWithPrice := range gasStationWithPricesMap {
			//Type available and target price reached?
			if (myGasStationWithPrice.gasPrice.Diesel > 10 && appUser.TargetDiesel >= myGasStationWithPrice.gasPrice.Diesel) ||
				(myGasStationWithPrice.gasPrice.E10 > 10 && appUser.TargetE10 >= myGasStationWithPrice.gasPrice.E10) ||
				(myGasStationWithPrice.gasPrice.E5 > 10 && appUser.TargetE5 >= myGasStationWithPrice.gasPrice.E5) {
				//Distance match?
				distance, _ := myGasStationWithPrice.gasStation.CalcDistanceBearing(appUser.Latitude, appUser.Longitude)
				if appUser.SearchRadius >= distance {
					//log.Printf("Match found: %v %v %v %v %v %v\n", distance, myGasStationWithPrice.gasStation.Brand, myGasStationWithPrice.gasStation.Place,
					//	myGasStationWithPrice.gasPrice.Diesel, myGasStationWithPrice.gasPrice.E10, myGasStationWithPrice.gasPrice.E5)
					gsMatches = append(gsMatches, myGasStationWithPrice)
				}
			}
		}
		if len(gsMatches) > 0 {
			myTitle := "Gas price matches found."
			myMessage := ""
			myDatas := []NotificationData{}
			for _, gsMatch := range gsMatches {
				myMessage = myMessage + fmt.Sprintf("Location: %v, E5: %v, E10: %v, Diesel: %v\n", gsMatch.gasStation.Place, (float64(gsMatch.gasPrice.E5)/1000), (float64(gsMatch.gasPrice.E10)/1000), (float64(gsMatch.gasPrice.Diesel)/1000))
				myNotificationData := NotificationData{GasStationID: gsMatch.gasStation.ID, StationName: gsMatch.gasStation.StationName, Brand: gsMatch.gasStation.Brand,
					Street: gsMatch.gasStation.Street, Place: gsMatch.gasStation.Place, HouseNumber: gsMatch.gasStation.HouseNumber, PostCode: gsMatch.gasStation.PostCode,
					Latitude: gsMatch.gasStation.Latitude, Longitude: gsMatch.gasStation.Longitude, Timestamp: time.Now(), E5: gsMatch.gasPrice.E5, E10: gsMatch.gasPrice.E10,
					Diesel: gsMatch.gasPrice.Diesel}
				myDatas = append(myDatas, myNotificationData)
			}
			myDataJson, err := json.Marshal(myDatas)
			if err != nil {
				log.Printf("Json marshal failed: %v", err)
				myDataJson = []byte("[]")
			}
			myNotificationMsg := NotificationMsg{UserUuid: appUser.Uuid, Message: myMessage, Title: myTitle, DataJson: string(myDataJson)}
			if len(string(myDataJson)) > 3500 {
				log.Printf("App User id: %v, Matches: %v, Json length: %v\n", appUser.Uuid, len(gsMatches), len(string(myDataJson)))
			}
			myNotificationMsgs = append(myNotificationMsgs, myNotificationMsg)
		}
	}
	StoreNotifications(myNotificationMsgs)
}
