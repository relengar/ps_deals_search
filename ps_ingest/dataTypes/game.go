package datatypes

import "time"

type Game struct {
	Name          string    `json:"name"`
	Price         float32   `json:"price"`
	Currency      string    `json:"currency"`
	OriginalPrice float32   `json:"originalPrice"`
	URL           string    `json:"url"`
	Description   string    `json:"description"`
	Rating        float32   `json:"rating"`
	RatingsSum    int       `json:"ratingsNum"`
	Expiration    time.Time `json:"expiration"`
	Platforms     []string  `json:"platforms"`
}
