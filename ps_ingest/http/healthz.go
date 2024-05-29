package http

import (
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

func healthzHanlder(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Ok")
}

func StartHealthz() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHanlder)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if errors.Is(err, net.ErrClosed) {
		log.Info().Msg("Server closed")
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		return err
	}

	return nil
}
