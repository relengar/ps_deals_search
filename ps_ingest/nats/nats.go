package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Connect() error
	Subscribe(string, nats.Handler) error
	Request(subject string, payload any, resp any) error
	Close()
}

type client struct {
	url   string
	token string
	conn  *nats.EncodedConn
	sub   *nats.Subscription
}

func (c *client) Connect() error {
	conn, err := nats.Connect(c.url, nats.Token(c.token))
	if err != nil {
		log.Error().Err(err).Str("url", c.url).Msg("Failed to connect to nats")
		return err
	}

	log.Info().Msg("Connected to nats")

	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Error().Err(err).Str("url", c.url).Msg("Failed to create encoded connection")
		return err
	}

	c.conn = nc
	return nil
}

func (c *client) Subscribe(subject string, handler nats.Handler) error {
	if c.conn == nil {
		err := fmt.Errorf("Can' subscribe to before connection is established")
		log.Error().Err(err).Str("subject", subject).Str("url", c.url).Msg("Failed to subscribe to subject")
		return err
	}

	sub, err := c.conn.Subscribe(subject, handler)
	if err != nil {
		log.Error().Err(err).Str("subject", subject).Str("url", c.url).Msg("Failed to subscribe to subject")
		return err
	}

	log.Info().Str("subject", subject).Msg("Subscribed")
	c.sub = sub

	return nil
}

func (c *client) Request(subject string, payload any, resp any) error {
	if c.conn == nil {
		err := fmt.Errorf("Can' send request to before connection is established")
		log.Error().Err(err).Str("subject", subject).Str("url", c.url).Msg("Failed nats request")
		return err
	}

	err := c.conn.Request(subject, payload, resp, 10*time.Second)
	if err != nil {
		return err
	}

	log.Debug().Msg("Got embedder response")
	return err
}

func (c *client) unsubscribe() error {
	if c.sub == nil {
		return nil
	}

	err := c.sub.Drain()
	if err != nil {
		return fmt.Errorf("Failed to process messages before unsubscribing %w", err)
	}

	return nil
}

func (c *client) Close() {
	if err := c.unsubscribe(); err != nil {
		log.Error().Err(err).Str("subject", c.sub.Subject).Str("url", c.url).Msg("Failed to unsbscribe from subject")
	}

	c.conn.Close()
}

type ConnectionConfig struct {
	Url   string
	Token string
}

func CreateClient(cfg ConnectionConfig) Client {
	return &client{cfg.Url, cfg.Token, nil, nil}
}
