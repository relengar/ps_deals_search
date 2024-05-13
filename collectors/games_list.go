package collectors

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

type gameCollectorConstructor func(string, chan Game) Crawler

type gameCrawler struct {
	collector           CollyCollector
	domain              string
	out                 chan Game
	createGameCollector gameCollectorConstructor
}

var hadGames = true

func (g *gameCrawler) Visit(path string) {
	defer close(g.out)

	visits := 0
	var err error

	for hadGames && visits < 3 {
		hadGames = false
		path, err = g.tick(path)
		if err != nil {
			break
		}
		visits += 1
	}

}

func (g *gameCrawler) tick(path string) (string, error) {
	log.Info().Str("path", path).Msg("Visiting deals page")
	url := g.toUrl(path)
	visitLog := log.Info().Str("url", url)
	visitLog.Msg("Visiting page")
	g.collector.Visit(url)
	g.collector.Wait()

	visitLog.Bool("had", hadGames).Str("path", path).Msg("Visited page")

	nextpath, err := getNextPath(path)
	if err != nil {
		return "", err
	}

	return nextpath, nil
}

func (g *gameCrawler) toUrl(path string) string {
	filters := "?FULL_GAME=storeDisplayClassification&PS4=targetPlatforms"
	return "https://" + g.domain + path + filters
}

func (g *gameCrawler) onHtml(h *colly.HTMLElement) {
	log.Info().Str("game", h.Text).Msg("I see a game")
	hadGames = true

	gameCollector := g.createGameCollector(g.domain, g.out)
	link := h.Attr("href")
	gameCollector.Visit("https://" + g.domain + link)
}

func getNextPath(path string) (string, error) {
	segments := strings.Split(path, "/")
	pageNum, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		log.Error().Str("path", path).Msg("Failed to parse page num from path")
		return "", err
	}

	nextNum := strconv.Itoa(pageNum + 1)
	segments = append(segments[:len(segments)-1], nextNum)
	nextPath := strings.Join(segments, "/")

	return nextPath, nil
}

func createGameCrawler(domain string, out chan Game) gameCrawler {
	c := colly.NewCollector(colly.AllowedDomains(domain), colly.Async(true))
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	createGameCollector := func(domain string, out chan Game) Crawler { return createGameCollector(domain, out) }
	collector := gameCrawler{c, domain, out, createGameCollector}

	c.OnHTML("main ul li a", collector.onHtml)

	c.OnRequest(func(r *colly.Request) {
		log.Info().Str("URL", r.URL.String()).Msg("Collecting games")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error().Err(err).Msg("Failed while crawling games list")
	})

	return collector
}
