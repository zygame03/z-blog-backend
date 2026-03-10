package stats

import "time"

type Config struct {
	SyncInterval time.Duration `mapstructure:"sync_interval" json:"sync_interval"`
}
