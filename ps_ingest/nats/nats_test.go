package nats

import (
	"encoding/json"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type response struct {
	Ok         bool        `json:"ok"`
	Embeddings [][]float64 `json:"embeddings"`
}

func TestNatsRequest(t *testing.T) {
	// GIVEN
	nc, client := setup(t)
	err := client.Connect()
	require.Nil(t, err)

	toSend := response{Ok: true, Embeddings: [][]float64{{0.123678, 0.3342, 0.6123673, 0.6312621}}}
	sub, err := nc.Subscribe("test", func(msg *nats.Msg) {
		payload, err := json.Marshal(toSend)
		assert.Nil(t, err)

		err = msg.Respond(payload)
		assert.Nil(t, err)
	})
	require.Nil(t, err)

	defer tearDown(t, nc, sub, client)

	// WHEN
	resp := response{}
	err = client.Request("test", nil, &resp)

	// THEN
	require.Nil(t, err)
	require.NotEmpty(t, resp)
	assert.True(t, resp.Ok)
	assert.Equal(t, resp.Embeddings, toSend.Embeddings)
}

func setup(t *testing.T) (*nats.Conn, Client) {
	conn, err := nats.Connect("nats://localhost:4222", nats.Token("token"))
	require.Nil(t, err)

	cfg := ConnectionConfig{Url: "nats://localhost:4222", Token: "token"}
	client := CreateClient(cfg)

	return conn, client
}

func tearDown(t *testing.T, conn *nats.Conn, sub *nats.Subscription, client Client) {
	if sub != nil {
		err := sub.Unsubscribe()
		assert.Nil(t, err)
	}
	if conn != nil {
		conn.Close()
	}
	if client != nil {
		client.Close()
	}
}
