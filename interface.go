package main

import "context"

type Cache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
}

//go:generate mockgen -destination=mocks/mock_result.go -source=interface.go -package=mocks main CacheResult
type CacheResult interface {
	GetResults(context.Context, string) (string, error)
}
