package main

import "context"

type Cache interface {
	Get(context.Context, string) (string, bool)
	Set(context.Context, string, string) error
}
