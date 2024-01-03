package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/subhasis/spotify-service/config"
	"github.com/subhasis/spotify-service/controllers"
	"github.com/subhasis/spotify-service/database"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// Load Configurations from config.json using Viper
	config.LoadAppConfig()

	// Initialize Database
	database.Connect(config.AppConfig.Database.ConnectionString)
	database.Migrate()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Register Routes
	RegisterProductRoutes(router)

	// Start the server
	log.Printf(fmt.Sprintf("Starting Server on port %d", config.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.Port), router))
}

func RegisterProductRoutes(router *mux.Router) {

	// Route for creating a track
	router.HandleFunc("/create-track/{isrc}", controllers.CreateTrackHandler).Methods("POST")
	// Route for retrieving metadata by ISRC
	router.HandleFunc("/get-metadata/{isrc}", controllers.GetMetadataByISRCHandler).Methods("GET")
	// Route for retrieving metadata by artist
	router.HandleFunc("/get-metadata-by-artist/{artist}", controllers.GetMetadataByArtistHandler).Methods("GET")

}
