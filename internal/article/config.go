package article

import (
	"sync/atomic"
	"time"
)

type ArticleConfig struct {
	cacheBaseTTL time.Duration
	cacheUserTTl time.Duration
}

var articleCfg atomic.Value

func init() {
	articleCfg.Store(defaultArticleConfig())
}

func defaultArticleConfig() ArticleConfig {
	return ArticleConfig{
		cacheBaseTTL: 1 * time.Minute,
		cacheUserTTl: 1 * time.Minute,
	}
}

func SetArticleConfig(cfg ArticleConfig) {
	articleCfg.Store(cfg)
}

func GetArticleConfig() ArticleConfig {
	return articleCfg.Load().(ArticleConfig)
}
