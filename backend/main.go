package main

import (
	"fmt"

	adapter_config "github.com/angular-and-go/pkd/adapter/config"
)

func init() {
	adapter_config.LoadEnvVariables()
}

func main() {
	fmt.Println("Hello World.")
}
