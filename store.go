package main

import (
	"context"

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
		store, err = stores.NewKafkaStore(ctx)
	case "file":
		loadEnv(fileEnvfile)
		store, err = stores.NewFileStore()
	}

	return store, err

}
