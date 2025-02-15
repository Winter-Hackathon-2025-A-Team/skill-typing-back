package repository

import (
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DbRepository {

	r := DbRepository{
		db: db,
	}
	return &r
}
