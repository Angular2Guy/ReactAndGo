/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
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
	"runtime"
	"syscall"
	"time"
)

var start = time.Now()

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
	log.Printf("Startup took: %vms\n", time.Since(start).Milliseconds())
	updateThreadPoolSize()
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

func updateThreadPoolSize() {
	if runtime.NumCPU() < 12 {
		runtime.GOMAXPROCS(12)
	}
}
