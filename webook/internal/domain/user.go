package domain

import "time"

type User struct {
	Id       int64
	Email    string
	PassWord string
	Nickname string
	Birthday time.Time
	AboutMe  string
	Phone    string
	Ctime    time.Time
}
