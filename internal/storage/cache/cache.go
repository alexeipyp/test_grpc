package httpcache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

func New(cacheLifetime time.Duration) (*bigcache.BigCache, error) {
	return bigcache.New(context.Background(), bigcache.DefaultConfig(cacheLifetime))
}
