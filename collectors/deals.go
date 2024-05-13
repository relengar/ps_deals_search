package collectors

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/gocolly/colly"
)

const dealsUrl = "/en-sk/pages/deals"

type CollyCollector interface {
	Visit(url string) error
	OnHTML(selector string, callback colly.HTMLCallback)
	OnRequest(callback colly.RequestCallback)
	Wait()
}

type Crawler interface {
	Visit(url string)
}

type DealCollector struct {
	mainCollector CollyCollector
	gameCrawler   Crawler
	startUrl      string
}

func (c *DealCollector) Start() {
	c.mainCollector.Visit(c.startUrl)
}

func (c *DealCollector) onHtml(h *colly.HTMLElement) {
	dealLink := h.Attr("href")

	log.Info().Str("Link", dealLink).Str("Name", h.Text).Msg("Checking deals category")

	c.gameCrawler.Visit(dealLink)
}

func CreateDealCollector(out chan Game) DealCollector {
	domain := os.Getenv("DOMAIN")
	c := colly.NewCollector(colly.AllowedDomains(domain))

	gameCrawler := createGameCrawler(domain, out)

	dealCollector := DealCollector{mainCollector: c, gameCrawler: &gameCrawler, startUrl: "https://" + domain + dealsUrl}

	selector := "main header a.psw-solid-link"
	c.OnHTML(selector, dealCollector.onHtml)

	c.OnRequest(func(r *colly.Request) { fmt.Println("Checking deals") })

	return dealCollector
}
