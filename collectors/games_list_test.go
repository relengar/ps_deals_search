package collectors

import (
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

func TestNextPath(t *testing.T) {
	validPath := "https://store.playstation.com/en-sk/category/322ea58b-1a57-4586-93e0-286d1e926394/2"
	nextPath := "https://store.playstation.com/en-sk/category/322ea58b-1a57-4586-93e0-286d1e926394/3"
	path, err := getNextPath(validPath)
	assert.Nil(t, err)
	assert.Equal(t, path, nextPath)

	invalidPath := "https://store.playstation.com/en-sk/category/322ea58b-1a57-4586-93e0-286d1e926394/2/unexpected_suffix"
	_, err = getNextPath(invalidPath)
	assert.NotNil(t, err)
}

type mockColly struct {
	visited bool
}

func (mc *mockColly) Visit(url string) error {
	mc.visited = true
	return nil
}

func (mc mockColly) OnHTML(selector string, callback colly.HTMLCallback) {
}

func (mc mockColly) OnRequest(callback colly.RequestCallback) {
}

func (mc mockColly) Wait() {
}

type mockCrawler struct {
	called bool
}

func (mc *mockCrawler) Visit(url string) {
	mc.called = true
}

func TestGamesListOnHTML(t *testing.T) {
	out := make(chan Game)

	mockColly := mockColly{}
	mc := mockCrawler{}
	gameCrawler := gameCrawler{&mockColly, "", out, func(s string, c chan Game) Crawler { return &mc }}

	hadGames = false

	h := &colly.HTMLElement{}
	gameCrawler.onHtml(h)

	assert.True(t, hadGames)
	assert.Equal(t, true, mc.called)
}

func TestGameListTick(t *testing.T) {
	out := make(chan Game)

	mockColly := mockColly{}
	mc := mockCrawler{}
	gameCrawler := gameCrawler{&mockColly, "", out, func(s string, c chan Game) Crawler { return &mc }}

	path := "https://store.playstation.com/en-sk/category/322ea58b-1a57-4586-93e0-286d1e926394/2"
	nextPath := "https://store.playstation.com/en-sk/category/322ea58b-1a57-4586-93e0-286d1e926394/3"

	next, _ := gameCrawler.tick(path)
	assert.Equal(t, next, nextPath)
	assert.True(t, mockColly.visited)
}
