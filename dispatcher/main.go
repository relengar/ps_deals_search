package dispatcher

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type queueClient interface {
	Publish(string, any) error
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

	err := d.queue.Publish(d.subject, msg)
	if err != nil {
		log.Error().Err(err).Any("body", msg).Any("msg", msg).Msg("Failed to send game message")
		return err
	}
	return nil
}

func (d *Dispatcher) Close() {
	d.queue.Close()
}

func CreateDispatcher(cfg DispatcherConfig) (Dispatcher, error) {
	conn, err := nats.Connect(cfg.Url, nats.Token(cfg.Token))

	if err != nil {
		log.Error().Err(err).Msg("Failed to establish nats connection")
		return Dispatcher{}, err
	}

	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Error().Err(err).Msg("Failed to establish encoded nats connection")
		return Dispatcher{}, err
	}

	return Dispatcher{nc, cfg.Subject}, nil
}
