package site

import (
	"my_web/backend/internal/http"
	"time"
)

type Config struct {
	CacheBaseTTL time.Duration `mapstructure:"cache_base_ttl"`
}

func Schema() *http.ModuleSchema {
	return &http.ModuleSchema{
		Name: "site",
		Fields: []*http.FieldSchema{
			http.NewNumberSchema(
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
