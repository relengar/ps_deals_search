package serialization

import (
	"ps_ingest/nats"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SerializerConfig struct {
	NatsClient         nats.Client
	NatsEncoderSubject string
}

type Serializer struct {
	natsClient         nats.Client
	natsEncoderSubject string
}

func (s *Serializer) Serialize(texts []string) {
	var encodings [][]float64

	err := s.natsClient.Request(s.natsEncoderSubject, texts, encodings)
	if err != nil {
		inputs := zerolog.Arr()
		for _, v := range texts {
			inputs.Str(v)
		}

		log.Error().Err(err).Array("texts", inputs).Msg("Failed to retrieve encodings")
	}

	// TODO: save to postgres
}

func CreateSerializer(cfg *SerializerConfig) Serializer {
	return Serializer{cfg.NatsClient, cfg.NatsEncoderSubject}
}
