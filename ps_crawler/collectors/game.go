package collectors

import (
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

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

		priceText := h.ChildText("span[data-qa=\"mfeCtaMain#offer0#finalPrice\"]")
		currency := getCurrency(priceText)
		price, err := getPrice(priceText)
		if err != nil {
			log.Error().Str("priceText", priceText).Str("game", name).Err(err).Msg("Failed to parse game priceText")
			return
		}

		origPriceText := h.ChildText("span[data-qa=\"mfeCtaMain#offer0#originalPrice\"]")
		originalPrice, err := getPrice(origPriceText)
		if err != nil {
			log.Error().Str("origPriceText", origPriceText).Str("game", name).Err(err).Msg("Failed to parse game origPriceText")
			return
		}

		ratingText := h.ChildText("span[data-qa=\"mfe-star-rating#overall-rating#average-rating\"]")
		rating, err := strconv.ParseFloat(ratingText, 32)
		if err != nil {
			log.Error().Str("ratingText", ratingText).Str("game", name).Err(err).Msg("Failed to parse game ratingText")
			return
		}

		ratingsSumText := h.ChildText("span[data-qa=\"mfe-star-rating#overall-rating#total-ratings\"]")
		ratingsSum, err := strconv.Atoi(strings.Split(h.ChildText("span[data-qa=\"mfe-star-rating#overall-rating#total-ratings\"]"), " ")[0])
		if err != nil {
			log.Error().Str("ratingsSumText", ratingsSumText).Str("game", name).Err(err).Msg("Failed to parse game ratingsSumText")
			return
		}

		description := h.ChildText("p[data-qa=\"mfe-game-overview#description\"]")
		expiration, err := toExpiration(h.ChildText("span[data-qa=\"mfeCtaMain#offer0#discountDescriptor\"]"))
		if err != nil {
			log.Error().Time("expiration", expiration).Str("game", name).Err(err).Msg("Failed to parse game expiration time")
			return
		}

		platformsText := h.ChildText("dd[data-qa=\"gameInfo#releaseInformation#platform-value\"]")
		platforms := strings.Split(platformsText, ",")

		game := Game{name, price, currency, originalPrice, url, description, float32(rating), ratingsSum, expiration, platforms}

		log.Info().Any("game", game).Msg("Collected game")
		out <- game
	})

	c.OnRequest(func(r *colly.Request) { url = r.URL.String() })

	return collector
}

func toExpiration(expirationText string) (time.Time, error) {
	dateText := strings.Replace(expirationText, "Offer ends ", "", 1)

	expiration, err := time.Parse("1/2/2006 15:04 MST", dateText)
	if err == nil {
		return expiration, nil
	}

	expiration, err = time.Parse("2/1/2006 15:04 MST", dateText)
	if err == nil {
		return expiration, nil
	}

	expiration, err = time.Parse("1/2/2006 15:04 PM MST", dateText)
	if err == nil {
		return expiration, nil
	}

	expiration, err = time.Parse("2/1/2006 15:04 PM MST", dateText)
	if err == nil {
		return expiration, nil
	}

	return time.Time{}, err
}

func getPrice(priceText string) (float32, error) {
	priceNumText := strings.Join(strings.Split(priceText, "")[1:], "")
	price, err := strconv.ParseFloat(priceNumText, 32)
	if err != nil {
		return float32(0), err
	}

	return float32(price), nil
}

func getCurrency(priceText string) string {
	return strings.Split(priceText, "")[0]
}
