package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_PARAMS")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database. ")
	}

	if DB == nil || DB.Error != nil {
		log.Fatalf("DB is %v", DB)
		//log.Fatal(DB.Error.Error())
	}

	//log.Default().Printf("%v", DB.Migrator().CurrentDatabase())
}
