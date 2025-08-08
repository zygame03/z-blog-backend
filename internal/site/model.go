package site

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
	Text      string    `json:"text"`
	OnlineAt  time.Time `josn:"online_at"`
	OfflineAt time.Time `json:"offline_at"`
}
