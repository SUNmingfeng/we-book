package dao

import "gorm.io/gorm"

// 这里dao与gorm强耦合了，后面要改为接口
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
