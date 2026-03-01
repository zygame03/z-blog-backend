package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(fpath, fname string, cfg any) error {
	v := viper.New()

	v.AddConfigPath(fpath)
	v.SetConfigName(fname)
	if err := v.ReadInConfig(); err != nil {
		log.Println("读取配置文件失败", fpath, fname, err)
		return err
	}

	if err := v.Unmarshal(cfg); err != nil {
		log.Println("解析配置文件失败", err)
		return err
	}

	log.Println("配置读取成功")
	return nil
}
