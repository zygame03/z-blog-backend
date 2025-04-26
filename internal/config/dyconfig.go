package config

import (
	"my_web/backend/internal/article"
	"my_web/backend/internal/websiteData"
	"sync/atomic"
)

type Dyconfig struct {
	Article  article.ArticleConfig         `mapstructure:"article"`
	SiteData websiteData.WebsiteDataConfig `mapstructure:"site_data"`
}

var config atomic.Value

func GetConfig() *Dyconfig {
	return config.Load().(*Dyconfig)
}

func LoadDyConfig(fpath, fname string) error {
	var cfg Dyconfig
	err := LoadConfig(fpath, fname, &cfg)
	if err != nil {
		return err
	}
	config.Store(&cfg)
	return nil
}

func SaveDyConfig() {
	// TODO
}
