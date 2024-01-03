package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/subhasis/spotify-service/config"
	"github.com/subhasis/spotify-service/entities"
)

func GetSpotifyTrackInfo(isrc string) (entities.Track, error) {

	const spotifyAPIURL = "https://api.spotify.com/v1/search"

	accessToken, err := getSpotifyAccessToken()
	if err != nil {
		return entities.Track{}, err
	}

	apiURL := fmt.Sprintf("%s?q=irsc:%s&type=track", spotifyAPIURL, url.QueryEscape(isrc))

	fmt.Println("REQUEST : ", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return entities.Track{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entities.Track{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entities.Track{}, fmt.Errorf("Spotify API request failed with status code %d", resp.StatusCode)
	}

	var spotifyResponse entities.SpotifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&spotifyResponse); err != nil {
		return entities.Track{}, err
	}

	if len(spotifyResponse.Tracks.Items) == 0 {
		return entities.Track{}, fmt.Errorf("no track found for title: %s", isrc)
	}

	highestPopularityIndex := findHighestPopularityIndex(spotifyResponse)
	if highestPopularityIndex < 0 {
		return entities.Track{}, fmt.Errorf("no track found with popularity")
	}

	track := entities.Track{
		ISRC:            isrc, //spotifyResponse.Tracks.Items[highestPopularityIndex].ExternalIDs.ISRC,
		SpotifyImageURI: "",   // TODO : which url to consider
		Title:           spotifyResponse.Tracks.Items[highestPopularityIndex].Name,
		Popularity:      spotifyResponse.Tracks.Items[highestPopularityIndex].Popularity,
		ArtistNameList:  extractArtistNames(spotifyResponse.Tracks.Items[highestPopularityIndex].Artists),
	}

	return track, nil
}

func getSpotifyAccessToken() (string, error) {

	resp, err := http.PostForm("https://accounts.spotify.com/api/token", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {config.AppConfig.Spotify.OAuth.ClientID},
		"client_secret": {config.AppConfig.Spotify.OAuth.ClientSecret},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Spotify authentication request failed with status code %d", resp.StatusCode)
	}

	// Decode the JSON response and extract the access token
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

// extractArtistNames extracts artist names from the API response
func extractArtistNames(artists []entities.Artist) entities.ArtistNameList {
	var names entities.ArtistNameList
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return names
}

func findHighestPopularityIndex(items entities.SpotifyResponse) int {
	if len(items.Tracks.Items) == 0 {
		return -1
	}

	highestPopularity := items.Tracks.Items[0].Popularity
	highestPopularityIndex := 0

	for i, item := range items.Tracks.Items {
		if item.Popularity > highestPopularity {
			highestPopularity = item.Popularity
			highestPopularityIndex = i
		}
	}

	return highestPopularityIndex
}
