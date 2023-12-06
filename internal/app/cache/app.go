package cacheapp

import "go.uber.org/zap"

type Cache interface {
	Close() error
}

type App struct {
	logger *zap.Logger
	cache  Cache
}

func New(cache Cache, logger *zap.Logger) *App {
	return &App{cache: cache, logger: logger}
}

func (a *App) Stop() {
	if err := a.cache.Close(); err != nil {
		a.logger.Panic("failed to stop cache", zap.Error(err))
	}
	a.logger.Info("cache successfully closed")
}
