package config

import (
	"fmt"
	"my_web/backend/internal/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(fpath, fname string, cfg any) error {
	v := viper.New()

	v.AddConfigPath(fpath)
	v.SetConfigName(fname)
	if err := v.ReadInConfig(); err != nil {
		logger.Error(
			"读取配置文件失败",
			zap.String("文件路径", fmt.Sprintf(fpath, fname)),
			zap.Error(err),
		)
		return err
	}

	if err := v.Unmarshal(cfg); err != nil {
		logger.Error(
			"解析配置文件失败",
			zap.Error(err),
		)
		return err
	}

	logger.Info("配置读取成功")
	return nil
}
