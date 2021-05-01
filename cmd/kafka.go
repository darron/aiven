package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/segmentio/kafka-go"
)

func Producer(cfg Config) (*kafka.Writer, error) {

	// Get the dialer with TLS config.
	dialer, err := Dialer()
	if err != nil {
		return nil, err
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{cfg.KafkaHostname},
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})
	return w, nil
}

func Consumer(cfg Config) (*kafka.Reader, error) {

	// Get the dialer with TLS config.
	dialer, err := Dialer()
	if err != nil {
		return nil, err
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaHostname},
		GroupID: cfg.KafkaConsumerGroup,
		Topic:   cfg.KafkaTopic,
		Dialer:  dialer,
	})

	return r, nil
}

func Dialer() (*kafka.Dialer, error) {
	var dialer kafka.Dialer

	// Load the certs from the filesystem.
	cert, err := tls.LoadX509KeyPair("service.cert", "service.key")
	if err != nil {
		return &dialer, err
	}

	// Let's load the CA
	ca, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		return &dialer, err
	}
	caCerts := x509.NewCertPool()
	caCerts.AppendCertsFromPEM(ca)

	dialer = kafka.Dialer{
		Timeout: 15 * time.Second,
		// DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCerts,
		},
	}
	return &dialer, nil
}
