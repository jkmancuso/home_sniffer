package main

func NewCache(_ string) Cache {
	return NewRedisCache()
}
