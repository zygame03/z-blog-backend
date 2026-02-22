package data

import (
	"sync/atomic"
	"time"
)

type SitedataConfig struct {
	cacheBaseTTL time.Duration
}

var config atomic.Value

func init() {
	config.Store(defaultConfig())
}

func defaultConfig() SitedataConfig {
	return SitedataConfig{
		cacheBaseTTL: 600 * time.Minute,
	}
}

func setConfig(cfg SitedataConfig) {
	config.Store(cfg)
}

func getConfig() SitedataConfig {
	return config.Load().(SitedataConfig)
}
