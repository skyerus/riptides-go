package models

type SpotifyCredentials struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyAuthorization struct {
	Code string `json:"code"`
}

type Play struct {
	URI string `json:"uri"`
}

type SpotifyPlay struct {
	URIs [1]string `json:"uris"`
}

type SpotifySearch struct {
	Tracks SpotifyTrack `json:"tracks"`
}

type SpotifyTrack struct {
	Items []SpotifyItem `json:"items"`
}

type SpotifyItem struct {
	Album SpotifyAlbum `json:"album"`
	Artists []SpotifyArtist `json:"artists"`
	Name string `json:"name"`
	DurationMs int `json:"duration_ms"`
	URI string `json:"uri"`
}

type SpotifyAlbum struct {
	Images []SpotifyImage `json:"images"`
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyImage struct {
	Height int `json:"height"`
	Width int `json:"width"`
	URL string `json:"url"`
}

type SpotifySearchSimple struct {
	Name string `json:"name"`
	Artist string `json:"artist"`
	URI string `json:"uri"`
	DurationMs int `json:"duration_ms"`
	Image SpotifyImage `json:"image"`
}
