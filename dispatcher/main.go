package dispatcher

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type queueClient interface {
	Publish(string, []byte) error
	Close()
}

type DispatcherConfig struct {
	Url     string
	Subject string
	Token   string
}

type Dispatcher struct {
	queue   queueClient
	subject string
}

func (d *Dispatcher) Dispatch(msg any) error {
	log.Info().Any("body", msg).Msg("Dispatching")

	value, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Any("body", msg).Msg("Failed to marshal game message")
		return err
	}

	err = d.queue.Publish(d.subject, value)
	if err != nil {
		log.Error().Err(err).Any("body", msg).Str("payload", string(value)).Msg("Failed to send game message")
		return err
	}
	return nil
}

func (d *Dispatcher) Close() {
	d.queue.Close()
}

func CreateDispatcher(cfg DispatcherConfig) (Dispatcher, error) {
	nc, err := nats.Connect(cfg.Url, nats.Token(cfg.Token))

	if err != nil {
		log.Error().Err(err).Msg("Failed to establish nats connection")
		return Dispatcher{}, err
	}
	return Dispatcher{nc, cfg.Subject}, nil
}
