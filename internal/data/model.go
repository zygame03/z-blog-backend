package data

import "time"

// k-v type model
type Sitedata struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
}

// announcement
type Announcement struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string
	IsDelete  bool
}
