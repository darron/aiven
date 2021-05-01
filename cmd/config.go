package cmd

import (
	"time"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	CommonConfig
	GatherConfig
	StoreConfig
}

type CommonConfig struct {
	KafkaTopic    string `env:"KAFKA_TOPIC,default=sites"`
	KafkaHostname string `env:"KAFKA_HOST,required"`
	KafkaCert     string `env:"KAFKA_CERT,default=service.cert"`
	KafkaKey      string `env:"KAFKA_KEY,default=service.key"`
	KafkaCA       string `env:"KAFKA_CA,default=ca.pem"`
}

type StoreConfig struct {
	KafkaConsumerGroup string `env:"KAFKA_CONSUMER_GROUP,default=storage"`
}

type GatherConfig struct {
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT,default=5s"`
	SitesList   string        `env:"SITES_LIST,default=websites.csv"`
}

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
