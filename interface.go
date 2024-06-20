package main

import "context"

type Cache interface {
	Get(context.Context, string) (ipInfo, bool)
	Set(context.Context, string, string) error
}
