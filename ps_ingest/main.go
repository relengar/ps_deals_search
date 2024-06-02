package main

import (
	"os"
	"os/signal"
	"ps_ingest/http"
	"ps_ingest/nats"
	"ps_ingest/serialization"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type config struct {
	QueueUrl     string `env:"QUEUE_URL,notEmpty"`
	QueueSubject string `env:"QUEUE_SUBJECT,notEmpty"`
	QueueToken   string `env:"QUEUE_TOKEN,notEmpty,unset"`
}

func main() {
	log.Info().Msg("Starting")

	cfg := loadConfig()

	// ctx, cancel := context.WithCancel(context.Background())
	stopChan := make(chan struct{}, 1)
	out := make(chan any)

	natsCfg := nats.ConnectionConfig{Url: cfg.QueueUrl, Token: cfg.QueueToken}
	natsClient := nats.CreateClient(natsCfg)
	err := natsClient.Connect()
	if err != nil {
		return
	}

	serializerCfg := serialization.SerializerConfig{NatsClient: natsClient, NatsEncoderSubject: cfg.QueueSubject}
	serializer := serialization.CreateSerializer(&serializerCfg)

	err = natsClient.Subscribe(cfg.QueueSubject, func(v any) {
		log.Info().Any("msg", v).Msg("Received message")
		serializer.Serialize([]string{"meh"})
	})
	if err != nil {
		natsClient.Close()
		return
	}

	// subCfg := subscriber.SubscriberConfig{Url: cfg.QueueUrl, Subject: cfg.QueueSubject, Token: cfg.QueueToken, Services: encoder}
	// sub := subscriber.CreateSubscriber(out, subCfg)
	// err := sub.Subscribe()
	// if err != nil {
	// 	return
	// }

	// TODO: implement connector and transformation to some persistance layer
	go func() {
		for v := range out {
			log.Info().Any("msg", v).Msg("Got message from subscriber")
		}
	}()

	go http.StartHealthz()

	go onTerminate(stopChan)

	<-stopChan
	// TODO: proper graceful exit
	natsClient.Close()
	log.Info().Msg("Exiting")

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

func onTerminate(stopChan chan struct{}) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	log.Info().Msg("Exiting process")
	stopChan <- struct{}{}
}
