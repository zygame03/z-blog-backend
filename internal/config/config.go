package config

import (
	"log"
	"my_web/backend/internal/infra"

	"github.com/spf13/viper"
)

// StaticConfig 应用配置结构
type StaticConfig struct {
	Httpserver infra.HttpserverConfig `mapstructure:"httpserver"`
	Database   infra.DatabaseConfig   `mapstructure:"database"`
	Redis      infra.RedisConfig      `mapstructure:"redis"`
}

// ReadConfig 读取配置文件
func ReadConfig(fpath, fname, ftype string) (*StaticConfig, error) {
	viper.Reset()

	viper.AddConfigPath(fpath)
	viper.SetConfigName(fname)
	viper.SetConfigType(ftype)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("读取配置文件失败", fpath, fname, ftype, err)
		return nil, err
	}

	var cfg StaticConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Println("解析配置文件失败", err)
		return nil, err
	}

	log.Println("配置读取成功")
	return &cfg, nil
}
