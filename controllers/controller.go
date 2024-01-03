package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/subhasis/spotify-service/spotify"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/subhasis/spotify-service/database"
	"github.com/subhasis/spotify-service/entities"
)

func CreateTrackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isrc := vars["isrc"]

	track, err := spotify.GetSpotifyTrackInfo(isrc)
	if err != nil {
		sendErrorResponse(w, err, http.StatusBadRequest, "Error fetching spotify metadata")
		return
	}

	// Store in the database
	newTrack := formatQueryData(track)

	// TODO :  to be removed
	formattedItem, _ := json.MarshalIndent(newTrack, "", "    ")
	fmt.Println("newTrack ::  ", string(formattedItem))

	tx := addToDB(newTrack)
	if tx.Error != nil {
		sendErrorResponse(w, tx.Error, http.StatusInternalServerError, "Record cannot be inserted")
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"isrc": isrc, "msg": "created"}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)

}

func formatQueryData(track entities.Track) entities.Track {
	return entities.Track{
		ISRC:            track.ISRC,
		SpotifyImageURI: track.SpotifyImageURI,
		Title:           track.Title,
		ArtistNameList:  track.ArtistNameList,
		Popularity:      track.Popularity,
	}
}

func addToDB(newTrack entities.Track) (tx *gorm.DB) {
	return database.Instance.Create(&newTrack)
}

func GetMetadataByISRCHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isrc := vars["isrc"]

	var track entities.Track
	tx := database.Instance.Where("ISRC = ?", isrc).First(&track)
	if tx.Error != nil {
		sendErrorResponse(w, tx.Error, http.StatusBadRequest, "No records found")
		return
	}

	response := map[string]interface{}{
		"ISRC":            track.ISRC,
		"SpotifyImageURI": track.SpotifyImageURI,
		"Title":           track.Title,
		"ArtistNameList":  track.ArtistNameList,
		"Popularity":      track.Popularity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMetadataByArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artist := vars["artist"]

	var tracks []entities.Track
	tx := database.Instance.Where("artist_name_list LIKE ?", "%"+artist+"%").Find(&tracks)
	if tx.Error != nil {
		sendErrorResponse(w, tx.Error, http.StatusBadRequest, "No records found")
		return
	}

	response := make([]map[string]interface{}, len(tracks))
	for i, track := range tracks {
		response[i] = map[string]interface{}{
			"ISRC":            track.ISRC,
			"SpotifyImageURI": track.SpotifyImageURI,
			"Title":           track.Title,
			"ArtistNameList":  track.ArtistNameList,
			"Popularity":      track.Popularity,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, err error, statusCode int, msg string) {

	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "rejected",
		"error":  err.Error(),
		"msg":    msg,
	}

	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	w.Write([]byte(responseJSON))
}
