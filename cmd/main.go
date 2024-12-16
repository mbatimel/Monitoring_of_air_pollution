package main

import (
	"Monitoring_of_air_pollution/internal/storage/postgres"
	"context"
	"net/http"
	"time"

	"os"
	"os/signal"
	"syscall"

	"Monitoring_of_air_pollution/internal/config"
	"Monitoring_of_air_pollution/internal/service"

	"github.com/rs/zerolog/log"
)

const serviceName = "weather_data"

func main() {

	log.Logger = config.Values().Logger().With().Str("serviceName", serviceName).Logger()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	postgresStorage, err := postgres.New(config.Values().Postgres, log.Logger)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	srv, err := service.NewServerConfig(config.Values().ServiceBind, postgresStorage)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to create server")
	}

	// Run server in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatal().Err(err).Msg("Server run failed")
		}
	}()

	log.Info().Msg("Server started")

	// Wait for interrupt signal to gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := srv.Close(); err != nil {
		log.Logger.Fatal().Err(err).Msg("Server shutdown failed")
	}

	select {
	case <-ctxShutDown.Done():
		if ctxShutDown.Err() == context.DeadlineExceeded {
			log.Info().Msg("Shutdown timed out")
		}
	default:
		log.Info().Msg("Server exited properly")
	}
}