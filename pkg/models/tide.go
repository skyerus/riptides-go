package models

type Tide struct {
	ID int `json:"id"`
	Name string `json:"name"`
	DateCreated string `json:"dateCreated"`
	User User `json:"user"`
	FavoritedBy []User `json:"favoritedBy"`
	Genres []int `json:"genres"`
	Tags []string `json:"tags"`
	Participants []User `json:"participants"`
	About string `json:"about"`
}

type Genre struct {
	ID int
	Name string `json:"name"`
}

type Tag struct {
	ID int
	Name string `json:"name"`
}
