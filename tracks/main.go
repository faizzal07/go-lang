package main

import (
	"barber/tracks/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/getInfo", getInfoHandler)
	http.ListenAndServe(":8080", nil)
}

func getInfoHandler(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		http.Error(w, "Region is required", http.StatusBadRequest)
		return
	}

	lastfmAPIKey := "6518402777c30c0da5c1f6f75c090223"
	lastfmURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=geo.gettoptracks&country=%s&api_key=%s&format=json", region, lastfmAPIKey)

	resp, err := http.Get(lastfmURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var lastfmResponse models.LastfmResponse
	if err := json.Unmarshal(body, &lastfmResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(lastfmResponse.Tracks.Track) == 0 {
		http.Error(w, "No top tracks found for the specified region", http.StatusNotFound)
		return
	}

	topTrack := lastfmResponse.Tracks.Track[46]
	artistName := topTrack.Artist.Name
	trackName := topTrack.Name

	artistInfoURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=artist.getinfo&artist=%s&api_key=%s&format=json", url.QueryEscape(artistName), lastfmAPIKey)
	resp, err = http.Get(artistInfoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var artistInfo models.ArtistInfo
	if err := json.Unmarshal(body, &artistInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	artistBio := artistInfo.Artist.Bio
	artistImageURL := artistInfo.Artist.ImageLinks[0].URL

	fmt.Printf("Top Track in %s\n", region)
	fmt.Printf("Track: %s\n", trackName)
	fmt.Printf("Artist: %s\n", artistName)
	fmt.Printf("Artist Bio: %v\n", artistBio.Summary)
	fmt.Printf("Artist Image: %s\n", artistImageURL)
	fmt.Fprintf(w, "Success")
}
