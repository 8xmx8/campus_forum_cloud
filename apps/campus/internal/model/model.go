package model

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DAO struct {
	DB *gorm.DB
	CC *redis.UniversalClient
}

func New(db *gorm.DB, cc *redis.UniversalClient) *DAO {
	return &DAO{
		CC: cc,
		DB: db,
	}
}

func (d *DAO) InitTable() error {

	if err := d.DB.AutoMigrate(&User{}); err != nil {
		return err
	}

	return nil
}
