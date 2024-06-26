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
	case "file":
		store, err = NewFileStore()
	}

	return store, err

}
