package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jkmancuso/home_sniffer/stores"
)

const fileEnvfile = "./file.env"
const kafkaEnvfile = "./kafka.env"

func NewStore(ctx context.Context, outputType string) (stores.Sender, error) {
	var store stores.Sender
	var err error

	switch outputType {
	case "kafka":
		loadEnv(kafkaEnvfile)

		if os.Getenv("KAFKA_CLIENT_CERT") == "TRUE" {
			tlsConfig, err := NewTLSConfig(
				os.Getenv("KAFKA_CLIENT_CERT"),
				os.Getenv("KAFKA_CLIENT_KEY"),
				os.Getenv("KAFKA_SERVER_CERT"))

			if err != nil {
				log.Fatalf("Unable to create TLS config: %v", err)
			}

			store, err = stores.NewKafkaStore(ctx, tlsConfig)

			if err != nil {
				log.Fatalf("Unable to get tls enabled %v store: %v", outputType, err)
			}

		} else {
			store, err = stores.NewKafkaStore(ctx)

			if err != nil {
				log.Fatalf("Unable to get %v store: %v", outputType, err)
			}

		}

	case "file":
		loadEnv(fileEnvfile)
		store, err = stores.NewFileStore()
	}

	return store, err

}
