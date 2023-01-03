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
	database.DB.First(&myGasStation, id)
	return myGasStation
}
