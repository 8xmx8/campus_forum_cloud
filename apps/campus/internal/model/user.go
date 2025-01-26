package model

type User struct {
	Id       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
	Age      uint8  `gorm:"default:0"`
	Gender   uint8  `gorm:"default:0"` // 性别
}

func (d *DAO) Register(u *User) error {
	return d.DB.Create(u).Error
}
