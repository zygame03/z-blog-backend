package article

import (
	"my_web/backend/internal/response"
	"time"
)

type Config struct {
	SyncInterval time.Duration `mapstructure:"sync_interval"`
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}

func Schema() *response.ModuleSchema {
	return &response.ModuleSchema{
		Name: "article",
		Fields: []*response.FieldSchema{
			response.NewNumberSchema(
				"syncInterval",
				"sync interval",
				"",
				-1,
				0,
				1,
			),
			response.NewNumberSchema(
				"cacheBaseTTL",
				"cache base ttl",
				"",
				-1,
				0,
				1,
			),
		},
	}
}

func init() {
	response.Register(Schema())
}
