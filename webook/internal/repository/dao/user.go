package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
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
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicate uint16 = 1062 //唯一索引冲突错误
		if me.Number == duplicate {
			return ErrDuplicateEmail //出现唯一索引冲突错误时返回邮箱冲突，增加唯一索引后需要修改
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var u User
	//First方法找不到会返回错误，Find方法找不到err==nil
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx *gin.Context, userid int64) (User, error) {
	var u User
	//First方法找不到会返回错误，Find方法找不到err==nil
	err := dao.db.WithContext(ctx).Where("id=?", userid).First(&u).Error
	fmt.Println("=============找到的数据=============")
	fmt.Printf("nickname:%v", u.Nickname)
	return u, err
}

func (dao *UserDAO) UpdateById(ctx *gin.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.ID).Updates(
		map[string]interface{}{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"` //自增主键
	Email    string `gorm:"unique"`                   //唯一索引
	PassWord string
	Nickname string `gorm:"type=varchar(64)"`
	Birthday int64
	AboutMe  string

	//UTC 0的毫秒数，所有地方都使用UTC 0时区存储，只在前端展示时转换时区
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}