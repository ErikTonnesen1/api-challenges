package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ErikTonnesen1/api-challenges/internal/app"
	"github.com/ErikTonnesen1/api-challenges/internal/database"
)

func main() {
	var cfg app.Config

	flag.StringVar(&cfg.Port, "port", ":3000", "API server port")
	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("APICH_DSN"), "Api-challenges Database DSN")
	cfg.Db.MaxOpenConns = 25
	cfg.Db.MaxIdleConns = 25
	cfg.Db.MaxIdleTime = 15 * time.Minute
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.OpenDB(
		cfg.Db.Dsn,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("database connection established")

	defer db.Close()

	application := app.New(cfg, logger, database.New(db))
	engine := application.StartGin()

	if ginErr := application.Serve(engine); ginErr != nil {
		log.Fatalf("Error starting Gin server: %v", ginErr)
	}
}
