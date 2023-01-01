package main

import (
	"fmt"

	config "github.com/angular-and-go/pkd/config"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	fmt.Println("Hello World.")
}
