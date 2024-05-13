package collectors

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

type Game struct {
	Name          string `json:"name"`
	Price         string `json:"price"`
	OriginalPrice string `json:"originalPrice"`
	URL           string `json:"url"`
	Description   string `json:"description"`
	Rating        string `json:"rating"`
	NumOfRatings  string `json:"ratingsNum"`
	Expiration    string `json:"expiration"`
}

type gameCollector struct {
	collector CollyCollector
	out       chan Game
}

func (c gameCollector) Visit(path string) {
	c.collector.Visit(path)
	c.collector.Wait()
}

func createGameCollector(domain string, out chan Game) gameCollector {
	c := colly.NewCollector(colly.AllowedDomains(domain))
	collector := gameCollector{c, out}

	var url string

	c.OnHTML("main", func(h *colly.HTMLElement) {
		name := h.ChildText("div.pdp-game-title h1")
		price := h.ChildText("span[data-qa=\"mfeCtaMain#offer0#finalPrice\"]")
		originalPrice := h.ChildText("span[data-qa=\"mfeCtaMain#offer0#originalPrice\"]")
		ratingText := h.ChildText("span[data-qa=\"mfe-star-rating#overall-rating#average-rating\"]")
		numOfRatingsText := strings.Split(h.ChildText("span[data-qa=\"mfe-star-rating#overall-rating#total-ratings\"]"), " ")[0]
		description := h.ChildText("p[data-qa=\"mfe-game-overview#description\"]")
		expirationText := h.ChildText("span[data-qa=\"mfeCtaMain#offer0#discountDescriptor\"]")
		expiration, err := toExpiration(expirationText)
		if err != nil {
			log.Error().Str("expiration", expirationText).Str("game", name).Err(err).Msg("Failed to parse game expiration time")
			return
		}

		game := Game{name, price, originalPrice, url, description, ratingText, numOfRatingsText, expiration}

		log.Info().Any("game", game).Msg("Collected game")
		out <- game
	})

	c.OnRequest(func(r *colly.Request) { url = r.URL.String() })

	return collector
}

func toExpiration(expirationText string) (string, error) {
	dateText := strings.Replace(expirationText, "Offer ends ", "", 1)

	expiration, err := time.Parse("1/2/2006 15:04 MST", dateText)
	if err == nil {
		return expiration.Format(time.RFC3339), nil
	}

	expiration, err = time.Parse("2/1/2006 15:04 MST", dateText)
	if err == nil {
		return expiration.Format(time.RFC3339), nil
	}

	expiration, err = time.Parse("1/2/2006 15:04 PM MST", dateText)
	if err == nil {
		return expiration.Format(time.RFC3339), nil
	}

	expiration, err = time.Parse("2/1/2006 15:04 PM MST", dateText)
	if err == nil {
		return expiration.Format(time.RFC3339), nil
	}

	return "", err
}
