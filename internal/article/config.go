package article

import (
	"sync/atomic"
	"time"
)

type ArticleConfig struct {
	syncInterval time.Duration
	cacheBaseTTL time.Duration
	cacheUserTTl time.Duration
}

var articleCfg atomic.Value

func init() {
	articleCfg.Store(defaultArticleConfig())
}

func defaultArticleConfig() ArticleConfig {
	return ArticleConfig{
		syncInterval: 10 * time.Second,
		cacheBaseTTL: 10 * time.Second,
		cacheUserTTl: 10 * time.Second,
	}
}

func SetArticleConfig(cfg ArticleConfig) {
	articleCfg.Store(cfg)
}

func GetArticleConfig() ArticleConfig {
	return articleCfg.Load().(ArticleConfig)
}
