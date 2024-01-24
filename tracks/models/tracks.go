package models

type Track struct {
	Name   string `json:"name"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
}

type LastfmResponse struct {
	Tracks struct {
		Track []Track `json:"track"`
	} `json:"tracks"`
}

type ArtistInfo struct {
	Artist ArtistDetails `json:"artist"`
}

type ArtistDetails struct {
	Name       string      `json:"name"`
	Bio        ArtistBio   `json:"bio"`
	ImageLinks []ImageLink `json:"image"`
}

type ArtistBio struct {
	Summary string `json:"summary"`
}

type ImageLink struct {
	URL    string `json:"#text"`
	Source string `json:"size"`
}
