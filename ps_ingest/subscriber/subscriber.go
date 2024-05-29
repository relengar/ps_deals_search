package subscriber

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/nats-io/nats.go"
)

type payload struct {
	Name          string    `json:"name"`
	Price         float32   `json:"price"`
	Currency      string    `json:"currency"`
	OriginalPrice float32   `json:"originalPrice"`
	URL           string    `json:"url"`
	Description   string    `json:"description"`
	Rating        float32   `json:"rating"`
	RatingsSum    int       `json:"ratingsNum"`
	Expiration    time.Time `json:"expiration"`
}

type Subscriber struct {
	out      chan any
	cfg      SubscriberConfig
	nc       *nats.EncodedConn
	sub      *nats.Subscription
	services Services
}

type Services interface {
	CallEncoder(texts []string) [][]float64
}

func (s *Subscriber) Subscribe() error {
	conn, err := nats.Connect(s.cfg.Url, nats.Token(s.cfg.Token))
	if err != nil {
		log.Error().Err(err).Str("subject", s.cfg.Subject).Str("url", s.cfg.Url).Msg("Failed to connect to nats")
		return err
	}
	log.Info().Msg("Connected to nats")

	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Error().Err(err).Str("subject", s.cfg.Subject).Str("url", s.cfg.Url).Msg("Failed to create encoded connection")
		return err
	}

	log.Info().Msg("Established encoded connection")
	s.nc = nc

	sub, err := nc.Subscribe(s.cfg.Subject, s.processMsg)
	if err != nil {
		log.Error().Err(err).Str("subject", s.cfg.Subject).Str("url", s.cfg.Url).Msg("Failed to subscribe to subject")
		return err
	}

	log.Info().Str("subject", s.cfg.Subject).Msg("Subscribed")
	s.sub = sub

	return nil
}

func (s *Subscriber) Stop() error {
	err := s.sub.Drain()
	if err != nil {
		return fmt.Errorf("Failed to process messages before unsubscribing %w", err)
	}

	err = s.sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("Failed to unsubscribe from nats %w", err)
	}

	s.nc.Close()

	return nil
}

func (s *Subscriber) processMsg(msg *payload) {
	log.Info().Any("msg", msg).Msg("Received nats message")
	// TODO: dispatch for db serialization

	s.services.CallEncoder(strings.Split(msg.Description, ". "))

	s.out <- msg
}

type SubscriberConfig struct {
	Url      string
	Subject  string
	Token    string
	Services Services
}

func CreateSubscriber(out chan any, cfg SubscriberConfig) Subscriber {
	return Subscriber{out, cfg, nil, nil, cfg.Services}
}
