package main

import (
	"context"
	"io"
)

func NewStore(ctx context.Context, outputType string) (io.Writer, error) {
	var store io.Writer
	var err error

	switch outputType {
	case "kafka":
		store, err = NewKafkaStore(ctx)
	}

	return store, err

}
