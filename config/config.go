package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port     int            `json:"port"`
	Database DatabaseConfig `json:"database"`
	Spotify  SpotifyConfig  `json:"spotify"`
}

type DatabaseConfig struct {
	ConnectionString string `json:"connection-string"`
}

type SpotifyConfig struct {
	OAuth OAuthConfig `json:"oauth"`
}

type OAuthConfig struct {
	ClientID     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
}

var AppConfig *Config

func LoadAppConfig() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	err = json.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Fatal("Error unmarshalling config file:", err)
	}
}
