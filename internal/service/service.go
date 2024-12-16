package service

import (
	"Monitoring_of_air_pollution/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"time"
)

var ErrChannelClosed = errors.New("channel is closed")

type Server interface {
	Run(ctx context.Context) error
	Close() error
}
type server struct {
	srv    *http.Server
	db     *postgres.Storage
	stopCh chan struct{}
}

func (s *server) Run(ctx context.Context) error {

	// Update weather forecast
	s.stopCh = make(chan struct{})
	go s.startCheckAirPollutionProcess()

	ch := make(chan error, 1)
	defer close(ch)
	go func() {
		ch <- s.srv.ListenAndServe()
	}()
	select {
	case err, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		if err != nil {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
	case <-ctx.Done():
		close(s.stopCh)
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context faild: %w", err)
		}

	}
	return nil
}
func (s *server) Close() error {
	close(s.stopCh)
	return s.srv.Close()
}
func (s *server) startCheckAirPollutionProcess() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			log.Info().Msg("Stopping background weather update process")
			return
		case <-ticker.C:
			err := s.UpdateAirPollutiont()
			if err != nil {
				log.Logger.Fatal().Err(err).Msg("failed to update weather forecast")
			} else {
				log.Info().Msg("Weather forecast updated successfully")
			}
		}
	}
}
func (s *server) UpdateAirPollutiont() error {
	fmt.Println("I AM HERE")
	return nil
}
func NewServerConfig(cfg string, pg *postgres.Storage) (Server, error) {
	srv := http.Server{
		Addr: net.JoinHostPort("localhost", cfg),
	}
	sv := server{
		srv: &srv,
		db:  pg,
	}

	return &sv, nil

}
