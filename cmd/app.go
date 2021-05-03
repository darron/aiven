package cmd

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

// App holds all of the dependencies we use in the app.
type App struct {
	C       Config
	DB      *sqlx.DB
	HTTP    *http.Client
	KReader *kafka.Reader
	KWriter *kafka.Writer
}

// GetAppConfig returns an *App we use in Store and Gather.
func GetAppConfig(cfg Config) (*App, error) {
	var a App

	// Postgres
	if cfg.AppType == "store" {
		// Connect to Postgres
		db, err := DBConnect(cfg)
		if err != nil {
			return &a, fmt.Errorf("postgres problem: %w", err)
		}
		a.DB = db
	} else {
		a.DB = nil
	}

	// HTTP Client
	if cfg.AppType == "gather" {
		a.HTTP = &http.Client{}
	} else {
		a.HTTP = nil
	}

	// Kafka
	switch cfg.AppType {
	case "store":
		r, err := Consumer(cfg)
		if err != nil {
			return &a, fmt.Errorf("kafka problem: %w", err)
		}
		a.KReader = r
	case "gather":
		w, err := Producer(cfg)
		if err != nil {
			return &a, fmt.Errorf("kafka problem: %w", err)
		}
		a.KWriter = w
	}

	// Set the config.
	a.C = cfg

	return &a, nil
}
