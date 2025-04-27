package websiteData

import (
	"my_web/backend/internal/global"
	"time"
)

// k-v type model
type WebsiteData struct {
	global.BaseModel
	Key   string
	Value string
}

// announcement
type Announcement struct {
	global.BaseModel
	Text      string
	OnlineAt  time.Time
	OfflineAt time.Time
}
