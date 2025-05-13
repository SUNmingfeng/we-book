package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDAO_Insert(t *testing.T) {
	testcases := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		u    User
		want error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				//传sql正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlmock.NewResult(1, 1))
				return db
			},
			ctx: context.Background(),
			u: User{
				Nickname: "nick",
			},
			want: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				//传sql正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			u: User{
				Nickname: "nick",
			},
			want: ErrDuplicateEmail,
		},
		{
			name: "插入失败",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				//传sql正则表达式
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx: context.Background(),
			u: User{
				Nickname: "nick",
			},
			want: errors.New("数据库错误"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}),
				&gorm.Config{
					DisableAutomaticPing:   true,
					SkipDefaultTransaction: true,
				})
			assert.NoError(t, err)
			dao := NewUserDAO(db)
			err = dao.Insert(tc.ctx, tc.u)
			assert.Equal(t, tc.want, err)
		})
	}
}
