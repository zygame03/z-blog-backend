package stats

import "my_web/backend/internal/global"

// Digital statistics
type NumStats struct {
	global.BaseModel
	Key   string `mapstructure:"key" json:"key"`
	Value int64  `mapstructure:"value" json:"value"`
}

type VisitorLog struct {
}
