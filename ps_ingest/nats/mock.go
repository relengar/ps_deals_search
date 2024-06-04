package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type MockNats struct {
	M                  mock.Mock
	EmbeddingsResponse any
}

func (n *MockNats) Connect() error {
	args := n.M.Called()
	return args.Error(0)
}

func (n *MockNats) Subscribe(subject string, handler nats.Handler) error {
	args := n.M.Called(subject, handler)
	return args.Error(0)
}

func (n *MockNats) Request(subject string, payload any, resp any) error {
	args := n.M.Called(subject, payload, resp)

	fakeEmbeddingResponse, _ := json.Marshal(n.EmbeddingsResponse)
	json.Unmarshal(fakeEmbeddingResponse, resp)

	return args.Error(0)
}

func (n *MockNats) Close() {
	n.M.Called()
}
