package serialization

import (
	datatypes "ps_ingest/dataTypes"
	"ps_ingest/nats"
	"ps_ingest/postgres"

	"github.com/rs/zerolog/log"
)

type embeddingsResponse struct {
	Ok         bool        `json:"ok"`
	Embeddings [][]float64 `json:"embeddings"`
}

type SerializerConfig struct {
	NatsClient         nats.Client
	NatsEncoderSubject string
	PgClient           postgres.PgClient
}

type NatsClient interface {
	Request(subject string, payload any, resp any) error
}

type PgClient interface {
	InsertGame(datatypes.Game) (int, error)
	InsertGameEmbedding(int, []float64) error
}

type Serializer struct {
	natsClient         NatsClient
	natsEncoderSubject string
	pgClient           PgClient
}

func (s *Serializer) Serialize(game datatypes.Game) {
	errLog := log.Error().Any("game", game)
	resp := embeddingsResponse{}

	gameId, err := s.pgClient.InsertGame(game)
	if err != nil {
		errLog.Err(err).Msg("Failed to insert ps game to postgres")
		return
	}

	err = s.natsClient.Request(s.natsEncoderSubject, []string{game.Description}, &resp)
	if err != nil {
		errLog.Err(err).Msg("Failed to retrieve encodings")
		return
	}

	if !resp.Ok {
		errLog.Any("resp", resp).Msg("Embedding failed")
		return
	}

	descriptionEmbedding := resp.Embeddings[0] // We've sent game description for embedding as first parameter
	err = s.pgClient.InsertGameEmbedding(gameId, descriptionEmbedding)
	if err != nil {
		errLog.Err(err).Int("embedding_length", len(descriptionEmbedding)).Msg("Failed to insert game description embedding to db")
		return
	}

	log.Info().Any("game", game).Int("id", gameId).Msg("Finished processing game")
}

func CreateSerializer(cfg *SerializerConfig) Serializer {
	return Serializer{cfg.NatsClient, cfg.NatsEncoderSubject, cfg.PgClient}
}
