package postgres

import (
	datatypes "ps_ingest/dataTypes"

	"github.com/stretchr/testify/mock"
)

type MockPgClient struct {
	M mock.Mock
}

func (pg *MockPgClient) Connect() error {
	args := pg.M.Called()
	return args.Error(0)
}

func (pg *MockPgClient) Close() {
	pg.M.Called()
}

func (pg *MockPgClient) InsertGame(game datatypes.Game) (int, error) {
	args := pg.M.Called(game)
	return args.Int(0), args.Error(1)
}

func (pg *MockPgClient) InsertGameEmbedding(gameId int, embeddings []float64) error {
	args := pg.M.Called(gameId, embeddings)
	return args.Error(0)
}
