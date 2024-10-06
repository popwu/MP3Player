package handler

import (
	"mp3-player/internal/libs"
)

type Handlers struct {
	storage   *libs.FileStorage
	cache     *libs.Cache
	cacheSize int64
}

func NewHandlers(storage *libs.FileStorage, cache *libs.Cache, cacheSize int64) *Handlers {
	return &Handlers{storage: storage, cache: cache, cacheSize: cacheSize}
}
