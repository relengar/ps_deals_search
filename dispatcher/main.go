package dispatcher

import "github.com/rs/zerolog/log"

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(msg any) error {
	log.Info().Any("body", msg).Msg("Dispatching")
	return nil
}

func CreateDispatcher() Dispatcher {
	return Dispatcher{}
}
