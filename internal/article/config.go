package article

import (
	"my_web/backend/internal/response"
	"time"
)

type ArticleConfig struct {
	SyncInterval time.Duration `mapstructure:"sync_interval"`
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}

func ArticleSchema() *response.ModuleSchema {
	return &response.ModuleSchema{
		Name: "article",
		Fields: []*response.FieldSchema{
			response.NewNumberSchema(
				"syncInterval",
				"syncInterval",
				"",
				-1,
				0,
				1,
			),
			response.NewNumberSchema(
				"cacheBaseTTL",
				"cacheBaseTTL",
				"",
				-1,
				0,
				1,
			),
		},
	}
}

func init() {
	response.Register(ArticleSchema())
}
