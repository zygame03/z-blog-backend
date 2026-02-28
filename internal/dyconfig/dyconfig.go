package dyconfig

import (
	"my_web/backend/internal/article"
	"sync/atomic"
)

type Dyconfig struct {
	Article article.ArticleConfig
}

var config atomic.Value

func init() {
	cfg, err := LocalLoadConfig("")
	if err != nil {
		// 重试 + 默认配置保底
	}
	config.Store(cfg)
}

func GetConfig() *Dyconfig {
	return config.Load().(*Dyconfig)
}

func LocalLoadConfig(path string) (*Dyconfig, error) {
	return &Dyconfig{}, nil
}

func LocalSaveConfig() {

}
