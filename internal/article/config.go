package article

import (
	"my_web/backend/internal/config"
	"time"
)

type ArticleConfig struct {
	syncInterval time.Duration
	cacheBaseTTL time.Duration
}

func ArticleSchema() *config.ModuleSchema {
	return &config.ModuleSchema{
		Name: "article",
		Fields: []*config.FieldSchema{
			config.NewNumberSchema(
				"syncInterval",
				"syncInterval",
				"",
				-1,
				0,
				1,
			),
			config.NewNumberSchema(
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
	config.Register(ArticleSchema())
}
