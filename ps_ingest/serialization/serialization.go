package serialization

import (
	datatypes "ps_ingest/dataTypes"
	"ps_ingest/nats"
	"ps_ingest/postgres"
	"strings"

	"github.com/rs/zerolog/log"
)

type embeddingsResponse struct {
	Ok         bool        `json:"ok"`
	Embeddings [][]float64 `json:"embeddings"`
}

type embeddingPayload struct {
	texts []string
}

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
	errLog := log.Error().Any("game", game)
	resp := embeddingsResponse{}
	sentences := strings.Split(game.Description, ".")

	err := s.natsClient.Request(s.natsEncoderSubject, sentences, &resp)
	if err != nil {
		errLog.Err(err).Msg("Failed to retrieve encodings")
		return
	}

	if !resp.Ok {
		errLog.Any("resp", resp).Msg("Embedding failed")
		return
	}

	err = s.pgClient.InsertGame(game, resp.Embeddings)
	if err != nil {
		errLog.Err(err).Int("embeddings_length", len(resp.Embeddings)).Msg("Failed to insert ps game to postgres")
	}
}

func CreateSerializer(cfg *SerializerConfig) Serializer {
	return Serializer{cfg.NatsClient, cfg.NatsEncoderSubject, cfg.PgClient}
}
