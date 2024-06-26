package main

import (
	"context"
	"io"
)

func NewStore(ctx context.Context, outputType string) io.Writer {
	switch outputType {
	case "kafka":
		return NewKafkaStore(ctx)
	default:
		return NewKafkaStore(ctx)
	}

}
