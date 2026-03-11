package user

import "time"

type Config struct {
	TokenTTL time.Duration `mapstructure:"token_ttl" json:"token_ttl"`
}
