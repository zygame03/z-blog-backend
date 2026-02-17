package data

import "time"

type Sitedata struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
}
