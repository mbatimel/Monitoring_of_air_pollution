package service

import (
	"Monitoring_of_air_pollution/internal/model"
	"Monitoring_of_air_pollution/internal/storage/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	apiKey = "ТВОЙ АЙПИ ЖОПКА"
	lat    = "55.7069" // Широта Ленинский проспект 6
	lon    = "37.5833" // Долгота Ленинский проспект 6
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
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/air_pollution?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Ошибка запроса")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: получен статус %d\n", resp.StatusCode)
		return err
	}

	var airPollution model.AirPollutionResponse
	err = json.NewDecoder(resp.Body).Decode(&airPollution)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Ошибка парсинга JSON")
		return err
	}

	fmt.Printf("Координаты: широта %.4f, долгота %.4f\n", airPollution.Coord.Lat, airPollution.Coord.Lon)
	for _, item := range airPollution.List {

		fmt.Printf("AQI: %d, CO: %.2f, PM2.5: %.2f, PM10: %.2f\n", item.Main.AQI, item.Components.CO, item.Components.PM2_5, item.Components.PM10)
		fmt.Printf("Дата: %s\n", time.Unix(item.Dt, 0).Format("2006-01-02 15:04:05"))
	}
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
