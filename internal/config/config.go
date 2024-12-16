package config

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/seniorGolang/gokit/env"
)

var configuration *Config

const formatJSON = "json"

type Config struct {
	ServiceBind string `env:"BIND_ADDR" envDefault:":9000" useFromEnv:"-"`
	Postgres    Postgres
	LogLevel    string `env:"LOGGER_LEVEL" envDefault:"debug"`
	LogFormat   string `env:"LOGGER_FORMAT" envDefault:""`
}
type Postgres struct {
	MaxConn         int    `env:"POSTGRES_MAX_CONN" envDefault:"25"`
	MasterAddress   string `env:"POSTGRES_MASTER_ADDRESS"`
	ReplicaAddress  string `env:"POSTGRES_REPLICA_ADDRESS"`
	DBName          string `env:"POSTGRES_DB_NAME"`
	UserName        string `env:"POSTGRES_USER_NAME_RW"`
	Password        string `env:"POSTGRES_PASSWORD_RW"`
	UserNameRO      string `env:"POSTGRES_USER_NAME_RO"`
	PasswordRO      string `env:"POSTGRES_PASSWORD_RO"`
	MaxIdleLifetime string `env:"POSTGRES_MAX_IDLE_LIFETIME" envDefault:"30s"`
	MaxLifetime     string `env:"POSTGRES_MAX_LIFETIME" envDefault:"3m"`
	PrepareCacheCap int    `env:"POSTGRES_PREPARE_CACHE_CAP" envDefault:"128"`
	CacheDuration   string `env:"POSTGRES_CACHE_DURATION" envDefault:"12h"`
}

func Values() Config {
	return *internalConfig()
}
func internalConfig() *Config {

	if configuration == nil {

		configuration = &Config{}

		if err := env.Parse(configuration); err != nil {
			panic(err)
		}
	}
	return configuration
}
func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.LogLevel); err == nil {
		level = newLevel
	}
	var out io.Writer = os.Stdout
	if cfg.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Timestamp().Logger()
}
