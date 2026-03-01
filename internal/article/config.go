package article

import (
	"my_web/backend/internal/response"
	"time"
)

type ArticleConfig struct {
	SyncInterval time.Duration `mapstructure:"syncInterval"`
	CacheBaseTTL time.Duration `mapstructure:"cacheBaseTTL"`
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
