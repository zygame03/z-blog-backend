package sitedata

import "time"

type Sitedata struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	key       string
	value     string
}
