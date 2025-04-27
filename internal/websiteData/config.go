package websiteData

import (
	"my_web/backend/internal/response"
	"time"
)

type WebsiteDataConfig struct {
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}

func Schema() *response.ModuleSchema {
	return &response.ModuleSchema{
		Name: "website_data",
		Fields: []*response.FieldSchema{
			response.NewNumberSchema(
				"cache_base_ttl",
				"缓存基本过期时间",
				"设置缓存基本过期时间",
				-1,
				0,
				1,
			),
		},
	}
}
