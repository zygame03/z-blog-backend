package global

import (
	"log"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Cors     CorsConfig     `mapstructure:"cors"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// CorsConfig CORS 配置
type CorsConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
	ExposeHeaders    []string `mapstructure:"exposeHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Port     uint   `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Protocol int    `mapstructure:"protocol"`
}

// ReadConfig 读取配置文件
func ReadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		log.Println("读取配置文件失败", path, err)
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Println("解析配置文件失败", err)
		return nil, err
	}

	log.Println("配置读取成功")
	return &cfg, nil
}
