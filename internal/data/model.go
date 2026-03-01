package data

import (
	"my_web/backend/internal/global"
)

// k-v type model
type Sitedata struct {
	global.BaseModel
	Key   string
	Value string
}

// announcement
type Announcement struct {
	global.BaseModel
	Text     string
	IsDelete bool
}
