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
		syncInterval: 24 * 60 * time.Second,
		cacheBaseTTL: 5 * 60 * time.Second,
		cacheUserTTl: 24 * 60 * time.Second,
	}
}

func setConfig(cfg ArticleConfig) {
	articleCfg.Store(cfg)
}

func getConfig() ArticleConfig {
	return articleCfg.Load().(ArticleConfig)
}
