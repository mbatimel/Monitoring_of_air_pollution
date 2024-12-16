package postgres

import (
	"Monitoring_of_air_pollution/internal/config"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Storage struct {
	connectManager *manager
}

func New(cfg config.Postgres, logger zerolog.Logger) (*Storage, error) {
	connectManager, err := newManager(logger, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "newManager error")
	}

	return &Storage{
		connectManager: connectManager,
	}, nil
}
