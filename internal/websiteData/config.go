package websiteData

import (
	"time"
)

type WebsiteDataConfig struct {
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}
