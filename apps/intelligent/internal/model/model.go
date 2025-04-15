package model

import (
	"gorm.io/gorm"
)

type DAO struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *DAO {
	return &DAO{
		DB: db,
	}
}

func (d *DAO) InitDBTable() error {
	tables := []interface{}{
		&ChatComment{},
		&ChatArticle{},
	}
	return d.DB.AutoMigrate(tables...)
}
