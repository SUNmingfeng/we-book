package wire

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("host=localhost"))
	if err != nil {
		panic(err)
	}
	return db
}
