package stats

import "time"

type Config struct {
	SyncInterval time.Duration `mapstructure:"sync_interval" json:"sync_interval"`
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl" json:"cache_base_ttl"`
}
