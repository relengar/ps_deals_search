package serialization

import (
	datatypes "ps_ingest/dataTypes"
	"ps_ingest/nats"
	"ps_ingest/postgres"
	"strings"

	"github.com/rs/zerolog/log"
)

type SerializerConfig struct {
	NatsClient         nats.Client
	NatsEncoderSubject string
	PgClient           postgres.PgClient
}

type Serializer struct {
	natsClient         nats.Client
	natsEncoderSubject string
	pgClient           postgres.PgClient
}

func (s *Serializer) Serialize(game datatypes.Game) {
	var embeddings [][]float64
	sentences := strings.Split(game.Description, ".")

	err := s.natsClient.Request(s.natsEncoderSubject, sentences, embeddings)
	if err != nil {
		log.Error().Err(err).Any("game", game).Msg("Failed to retrieve encodings")
		return
	}

	err = s.pgClient.InsertGame(game, embeddings)
	if err != nil {
		log.Error().Err(err).Any("game", game).Any("embeddings", embeddings).Msg("Failed to insert ps game to postgres")
	}
}

func CreateSerializer(cfg *SerializerConfig) Serializer {
	return Serializer{cfg.NatsClient, cfg.NatsEncoderSubject, cfg.PgClient}
}
