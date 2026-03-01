package config

import (
	"my_web/backend/internal/infra"
)

// StaticConfig 应用配置结构
type StaticConfig struct {
	Httpserver infra.HttpserverConfig `mapstructure:"httpserver"`
	Database   infra.DatabaseConfig   `mapstructure:"database"`
	Redis      infra.RedisConfig      `mapstructure:"redis"`
}

// ReadConfig 读取配置文件
func LoadStConfig(fpath, fname string) (*StaticConfig, error) {
	var cfg StaticConfig
	err := LoadConfig(fpath, fname, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
