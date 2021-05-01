package cmd

import (
	"github.com/joeshaw/envdecode"
)

type Config struct {
	CommonConfig
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

	return config, nil
}
