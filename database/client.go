package database

import (
	"log"

	"github.com/subhasis/spotify-service/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var err error

func Connect(connectionString string) {
	// added ?parseTime=true to parse dates while fetching date formats from db
	// else will cause error
	Instance, err = gorm.Open(mysql.Open(connectionString+"?parseTime=true"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

func Migrate() {
	Instance.AutoMigrate(&entities.Track{})
	log.Println("Database Migration Completed...")
}
