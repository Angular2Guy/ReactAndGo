package main

import (
	"angular-and-go/pkd/config"
	"angular-and-go/pkd/controller"
	"angular-and-go/pkd/cron"
	"angular-and-go/pkd/database"
	"angular-and-go/pkd/messaging"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateDB()
	messaging.Start()
	cron.Start()
}

func main() {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go controller.Start()

	<-quit
	log.Println("Shutting down server...")

	messaging.Stop()
	time.Sleep(2 * time.Second)

	log.Println("Server exiting")
}
