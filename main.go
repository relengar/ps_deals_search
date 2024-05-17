package main

import (
	"ps_crawler/collectors"
	"ps_crawler/dispatcher"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file loaded")
	}

	log.Info().Msg("Starting")
	out := make(chan collectors.Game)

	dispatcher, err := dispatcher.CreateDispatcher()
	if err != nil {
		return
	}

	go func() {
		dealCollector := collectors.CreateDealCollector(out)
		dealCollector.Start()
	}()

	for v := range out {
		dispatcher.Dispatch(v)
	}

	dispatcher.Close()
	log.Info().Msg("Done")
}
