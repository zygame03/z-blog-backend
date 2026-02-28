package article

import (
	"my_web/backend/internal/response"
	"time"
)

type ArticleConfig struct {
	syncInterval time.Duration
	cacheBaseTTL time.Duration
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
