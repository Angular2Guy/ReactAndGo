package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"react-and-go/pkd/config"
	"react-and-go/pkd/controller"
	"react-and-go/pkd/cron"
	"react-and-go/pkd/database"
	"react-and-go/pkd/database/dbmigrate"
	"react-and-go/pkd/messaging"
	"syscall"
	"time"
)

//go:embed public
var embeddedFiles embed.FS

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	dbmigrate.MigrateDB()
	messaging.Start()
	cron.Start()
}

func main() {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go controller.Start(getPublicFolder(embeddedFiles))

	<-quit
	log.Println("Shutting down server...")

	messaging.Stop()
	time.Sleep(2 * time.Second)

	log.Println("Server exiting")
}

func getPublicFolder(myEmbeddedFiles embed.FS) fs.FS {
	result, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return result
}
