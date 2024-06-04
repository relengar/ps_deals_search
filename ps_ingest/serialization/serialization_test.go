package serialization

import (
	"errors"
	datatypes "ps_ingest/dataTypes"
	"ps_ingest/nats"
	"ps_ingest/postgres"
	"testing"

	"github.com/stretchr/testify/mock"
)

var mockEmbeddings = [][]float64{{0.12321321}}

func TestGameSerialization(t *testing.T) {
	gameId := 1
	game := datatypes.Game{
		Description: "description",
	}
	cases := []struct {
		embeddings   embeddingsResponse
		gameInsertOk bool
		natsOk       bool
	}{
		{
			embeddings:   embeddingsResponse{Ok: true, Embeddings: mockEmbeddings},
			gameInsertOk: true,
			natsOk:       true,
		},
		{
			embeddings:   embeddingsResponse{Ok: true, Embeddings: mockEmbeddings},
			gameInsertOk: false,
		},
		{
			embeddings: embeddingsResponse{Ok: false, Embeddings: mockEmbeddings},
			natsOk:     true,
		},
		{
			embeddings: embeddingsResponse{Ok: true, Embeddings: mockEmbeddings},
			natsOk:     false,
		},
	}

	for _, c := range cases {
		serializer, nats, pg := setup(c.embeddings)
		pg.M.On("InsertGame", game).Return(gameId, createErr("pg mock failure", c.gameInsertOk))
		if c.gameInsertOk {
			nats.M.On("Request", "subject", []string{game.Description}, &embeddingsResponse{}).Return(createErr("pg mock failure", c.natsOk))
		}
		if c.gameInsertOk && c.embeddings.Ok && c.natsOk {
			pg.M.On("InsertGameEmbedding", gameId, c.embeddings.Embeddings[0]).Return(nil)
		}

		serializer.Serialize(game)

		nats.M.AssertExpectations(t)
		pg.M.AssertExpectations(t)
	}
}

func setup(embeddingsResponse embeddingsResponse) (Serializer, *nats.MockNats, *postgres.MockPgClient) {
	nats := &nats.MockNats{M: mock.Mock{}, EmbeddingsResponse: embeddingsResponse}
	pg := &postgres.MockPgClient{}
	return Serializer{nats, "subject", pg}, nats, pg
}

func createErr(msg string, isNil bool) error {
	if isNil {
		return nil
	}
	return errors.New(msg)

}
