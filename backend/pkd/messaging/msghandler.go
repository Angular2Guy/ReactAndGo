package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PriceUpdates struct {
	Useconds     int64
	Diesel       json.Number
	E5           json.Number
	E10          json.Number
	Diesel_delta float64
	E5_delta     float64
	E10_delta    float64
}

var client mqtt.Client

var gasPriceMsgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Message: %s received on topic: %s size: %d\n", msg.Payload(), msg.Topic(), len(msg.Payload()))
	var priceUpdateRawMap map[string]json.RawMessage
	if err := json.Unmarshal(msg.Payload(), &priceUpdateRawMap); err != nil {
		log.Printf("Message: %s received on topic: %s size: %d\n", msg.Payload(), msg.Topic(), len(msg.Payload()))
		log.Fatalf("Unmarshal failed: %v\n", err.Error())
	}
	priceUpdateMap := make(map[string]PriceUpdates)
	for key, value := range priceUpdateRawMap {
		myPriceUpdates := PriceUpdates{Useconds: 0, Diesel: "0", E5: "0", E10: "0", Diesel_delta: 0, E5_delta: 0, E10_delta: 0}
		if err := json.Unmarshal(value, &myPriceUpdates); err != nil {
			log.Printf("PriceUpdate: %v\n", string(value))
			log.Printf("Unmarshal failed: %v\n", err)
		} else {
			myPriceUpdates.Diesel = json.Number(strings.TrimSpace(strings.ReplaceAll(myPriceUpdates.Diesel.String(), ".", "")))
			myPriceUpdates.E5 = json.Number(strings.TrimSpace(strings.ReplaceAll(myPriceUpdates.E5.String(), ".", "")))
			myPriceUpdates.E10 = json.Number(strings.TrimSpace(strings.ReplaceAll(myPriceUpdates.E10.String(), ".", "")))
			priceUpdateMap[key] = myPriceUpdates
		}
	}
	log.Default().Printf("PriceUpdateMap: %v", priceUpdateMap)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}

func Start() {
	msgServerUrl := os.Getenv("MSG_PARAMS")
	msgClientId := os.Getenv("MSG_CLIENT_ID")
	msgServerUser := os.Getenv("MSG_SERVER_USER")
	msgServerPwd := os.Getenv("MSG_SERVER_PWD")
	options := mqtt.NewClientOptions()
	options.AddBroker(msgServerUrl)
	options.SetClientID(msgClientId)
	options.SetUsername(msgServerUser)
	options.SetPassword(msgServerPwd)
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client = mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Printf("Connection failed: %v\n", token.Error())
	} else {
		log.Printf("Connected to: %v id: %v\n", msgServerUrl, msgClientId)
	}

	msgGasPriceTopic := os.Getenv("MSG_GAS_PRICE_TOPIC")
	token = client.Subscribe(msgGasPriceTopic, 1, gasPriceMsgHandler)
	if token.Wait() && token.Error() != nil {
		log.Printf("Topic subription to topic: %v failed: %v", msgGasPriceTopic, token.Error().Error())
	} else {
		log.Printf("Subscribed to topic %s\n", msgGasPriceTopic)
	}
}

func Stop() {
	client.Disconnect(1000)
}

func SendMsg(msg string) {
	msgGasPriceTopic := os.Getenv("MSG_GAS_PRICE_TOPIC")
	client.Publish(msgGasPriceTopic, 0, false, msg)
}
