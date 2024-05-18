package main

import (
	"ps_crawler/collectors"
	"ps_crawler/dispatcher"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type config struct {
	Domain       string `env:"DOMAIN,notEmpty"`
	MaxPages     int    `env:"MAX_PAGES" envDefault:"-1"`
	QueueUrl     string `env:"QUEUE_URL,notEmpty"`
	QueueSubject string `env:"QUEUE_SUBJECT,notEmpty"`
	QueueToken   string `env:"QUEUE_TOKEN,notEmpty,unset"`
}

func main() {
	log.Info().Msg("Starting")
	cfg := loadConfig()

	dispatcherCfg := dispatcher.DispatcherConfig{
		Url:     cfg.QueueUrl,
		Subject: cfg.QueueSubject,
		Token:   cfg.QueueToken,
	}
	dispatcher, err := dispatcher.CreateDispatcher(dispatcherCfg)
	if err != nil {
		return
	}

	out := make(chan collectors.Game)
	go func() {
		collectorCfg := collectors.CollectorConfig{Domain: cfg.Domain, MaxPages: cfg.MaxPages}
		dealCollector := collectors.CreateDealCollector(out, collectorCfg)
		dealCollector.Start()
	}()

	for game := range out {
		dispatcher.Dispatch(game)
	}

	dispatcher.Close()
	log.Info().Msg("Done")
}

func loadConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file loaded")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to create config from env")
	}

	return cfg
}
