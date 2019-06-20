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
