package cmd

import (
	"time"

	"github.com/joeshaw/envdecode"
)

// Config holds all of the configuration for the application.
type Config struct {
	CommonConfig
	GatherConfig
	StoreConfig
}

// CommonConfig holds configuration common to all modes of operation.
type CommonConfig struct {
	KafkaSSLEnable bool   `env:"KAFKA_SSL_ENABLE,default=true"`
	KafkaTopic     string `env:"KAFKA_TOPIC,default=sites"`
	KafkaHostname  string `env:"KAFKA_HOST,required"`
	KafkaCert      string `env:"KAFKA_CERT,default=service.cert"`
	KafkaKey       string `env:"KAFKA_KEY,default=service.key"`
	KafkaCA        string `env:"KAFKA_CA,default=ca.pem"`
}

// StoreConfig holds configuration only used by `bin/aiven store`
type StoreConfig struct {
	KafkaConsumerGroup string `env:"KAFKA_CONSUMER_GROUP,default=storage"`
	PostgresURL        string `env:"POSTGRES_URL,required"`
	PostgresCA         string `env:"POSTGRES_CA,default=postgres-ca.pem"`
}

// GatherConfig holds configuration only used by `bin/aiven gather`
type GatherConfig struct {
	GatherDelay time.Duration `env:"GATHER_DELAY,default=30s"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT,default=5s"`
	SitesList   string        `env:"SITES_LIST,default=websites.csv"`
}

// Load gets configuration from the environment and returns Config and error.
func Load(configType string) (Config, error) {
	var config Config
	var common CommonConfig

	// Check for common config we need.
	err := envdecode.StrictDecode(&common)
	if err != nil {
		return config, err
	}
	config.CommonConfig = common

	// Check for specialized config.
	if configType == "store" {
		var store StoreConfig
		err := envdecode.StrictDecode(&store)
		if err != nil {
			return config, err
		}
		config.StoreConfig = store
	}

	if configType == "gather" {
		var gather GatherConfig
		err := envdecode.StrictDecode(&gather)
		if err != nil {
			return config, err
		}
		config.GatherConfig = gather
	}

	return config, nil
}
