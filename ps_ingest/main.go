package main

import (
	"os"
	"os/signal"
	"ps_ingest/http"
	"ps_ingest/nats"
	"ps_ingest/postgres"
	"ps_ingest/serialization"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type config struct {
	QueueUrl         string `env:"QUEUE_URL,notEmpty"`
	QueueSubject     string `env:"QUEUE_SUBJECT,notEmpty"`
	QueueToken       string `env:"QUEUE_TOKEN,notEmpty,unset"`
	EmbeddingSubject string `env:"EMBEDDING_SUBJECT,notEmpty"`
	PgUser           string `env:"PG_USER,notEmpty"`
	PgPassword       string `env:"PG_PASSWORD,notEmpty,unset"`
	PgHost           string `env:"PG_HOST,notEmpty"`
	PgDatabase       string `env:"PG_DATABASE,notEmpty"`
	HTTPPort         int    `env:"PORT" envDefault:"8080"`
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

	pgConfig := postgres.PgConfig{User: cfg.PgUser, Password: cfg.PgPassword, Host: cfg.PgHost, Database: cfg.PgDatabase}
	pgClient := postgres.CreatePgClient(pgConfig)
	err = pgClient.Connect()
	if err != nil {
		natsClient.Close()
		return
	}

	serializerCfg := serialization.SerializerConfig{NatsClient: natsClient, NatsEncoderSubject: cfg.EmbeddingSubject, PgClient: pgClient}
	serializer := serialization.CreateSerializer(&serializerCfg)

	err = natsClient.Subscribe(cfg.QueueSubject, serializer.Serialize)
	if err != nil {
		pgClient.Close()
		natsClient.Close()
		return
	}

	go func() {
		for v := range out {
			log.Info().Any("msg", v).Msg("Got message from subscriber")
		}
	}()

	go http.StartHealthz(cfg.HTTPPort)

	go onTerminate(stopChan)

	<-stopChan
	// TODO: proper graceful exit
	pgClient.Close()
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
