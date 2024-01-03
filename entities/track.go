package entities

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type ArtistNameList []string

func (o *ArtistNameList) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(string(bytes), ",")
	return nil
}
func (o ArtistNameList) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}

// Track model for the database
type Track struct {
	gorm.Model
	ISRC            string
	SpotifyImageURI string
	Title           string
	ArtistNameList  ArtistNameList `gorm:"type:varchar(100)"`
	Popularity      int
}

// ExternalURLs struct represents external URLs in the Spotify response
type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

// Artist struct represents an artist in the Spotify response
type Artist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

// Image struct represents an image in the Spotify response
type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// Album struct represents an album in the Spotify response
type Album struct {
	AlbumType            string       `json:"album_type"`
	Artists              []Artist     `json:"artists"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalURLs         ExternalURLs `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Image      `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	TotalTracks          int          `json:"total_tracks"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
}

// ExternalIDs struct represents external IDs in the Spotify response
type ExternalIDs struct {
	ISRC string `json:"isrc"`
}

// SpotifyResponseTrack struct represents a track in the Spotify response
type SpotifyResponseTrack struct {
	Album            Album        `json:"album"`
	Artists          []Artist     `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int          `json:"disc_number"`
	DurationMS       int          `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIDs      ExternalIDs  `json:"external_ids"`
	ExternalURLs     ExternalURLs `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsLocal          bool         `json:"is_local"`
	Name             string       `json:"name"`
	Popularity       int          `json:"popularity"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int          `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}

// Tracks struct represents the tracks section in the Spotify response
type Tracks struct {
	Href  string                 `json:"href"`
	Items []SpotifyResponseTrack `json:"items"`
}

// SpotifyResponse struct represents the entire Spotify API response
type SpotifyResponse struct {
	Tracks Tracks `json:"tracks"`
}
