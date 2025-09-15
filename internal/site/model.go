package site

import (
	"my_web/backend/internal/global"
	"time"
)

type Intro struct {
	global.BaseModel
	Content string `json:"content"`
}

// announcement
type Announcement struct {
	global.BaseModel
	Text      string    `json:"text"`
	OnlineAt  time.Time `json:"online_at"`
	OfflineAt time.Time `json:"offline_at"`
}

type DanmakuStatus int

const (
	Pending DanmakuStatus = iota
	Approved
	Rejected
)

type Danmaku struct {
	global.BaseModel
	Content  string        `json:"content"`
	SenderID string        `json:"sender_id"`
	IP       string        `json:"ip"`
	Status   DanmakuStatus `json:"status"`
}
