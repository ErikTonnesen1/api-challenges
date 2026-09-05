package app

import (
	"log/slog"
	"time"

	"github.com/ErikTonnesen1/api-challenges/internal/database"
)

type Config struct {
	Port string
	Db   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}

type Application struct {
	Config Config
	Logger *slog.Logger
	Models database.Models
}

func New(cfg Config, logger *slog.Logger, models database.Models) *Application {
	return &Application{
		Config: cfg,
		Logger: logger,
		Models: models,
	}

}
