package home

import "gorm.io/gorm"

type database struct {
	*gorm.DB
}

func newDatabase(db *gorm.DB) *database {
	return &database{
		db,
	}
}
