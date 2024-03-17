package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//创建数据
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"` //自增主键
	Email    string `gorm:"unique"`                   //唯一索引
	PassWord string

	//UTC 0的毫秒数，所有地方都使用UTC 0时区存储，只在前端展示时转换时区
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}
