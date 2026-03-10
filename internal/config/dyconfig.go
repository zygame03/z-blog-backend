package config

import (
	"fmt"
	"my_web/backend/internal/article"
	"my_web/backend/internal/site"
	"my_web/backend/internal/stats"
	"my_web/backend/internal/user"
	"reflect"
	"sync/atomic"
)

type Dyconfig struct {
	Article  article.Config `mapstructure:"article"`
	SiteData site.Config    `mapstructure:"site_data"`
	Stats    stats.Config   `mapstructure:"stats"`
	User     user.Config    `mapstructure:"user"`
}

var config atomic.Value
var configMap = map[string]func() any{
	"article":   func() any { return &article.Config{} },
	"site_data": func() any { return &site.Config{} },
	"stats":     func() any { return &stats.Config{} },
}

func SetConfigWithModule(module string, conf any) error {
	oldConf := GetConfig()
	newConf := *oldConf
	v := reflect.ValueOf(&newConf).Elem()

	f := v.FieldByName(module)
	if !f.IsValid() {
		return fmt.Errorf("module not found")

	}
	if f.Type() != reflect.TypeOf(conf) {
		return fmt.Errorf("type mismatch")
	}
	f.Set(reflect.ValueOf(conf))
	config.Store(&newConf)
	return nil
}

func GetConfig() *Dyconfig {
	return config.Load().(*Dyconfig)
}

func GetArticleConfig() *article.Config {
	return &GetConfig().Article
}

func GetSiteDataConfig() *site.Config {
	return &GetConfig().SiteData
}

func GetStatsConfig() *stats.Config {
	return &GetConfig().Stats
}

func GetUserConfig() *user.Config {
	return &GetConfig().User
}

func NewConfigByModule(module string) any {
	if fn, ok := configMap[module]; ok {
		return fn()
	}
	return nil
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
