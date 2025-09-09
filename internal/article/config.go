package article

import (
	"my_web/backend/internal/http"
	"time"
)

type Config struct {
	SyncInterval time.Duration `mapstructure:"sync_interval"`
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}

func Schema() *http.ModuleSchema {
	return &http.ModuleSchema{
		Name: "article",
		Fields: []*http.FieldSchema{
			http.NewNumberSchema(
				"syncInterval",
				"sync interval",
				"",
				1000000000000,
				0,
				1,
			),
			http.NewNumberSchema(
				"cacheBaseTTL",
				"cache base ttl",
				"",
				1000000000000,
				0,
				1,
			),
		},
	}
}

func init() {
	http.Register(Schema())
}
