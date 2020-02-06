package models

import "time"

type ApplyMara struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	RealName  string    `form:"real_name" xorm:"not null default '' comment('姓名') VARCHAR(20)" validate:"required"`
	Code      string    `form:"code" xorm:"not null default '' comment('学号') VARCHAR(20)" validate:"required"`
	ClassName string    `form:"class_name" xorm:"not null default '' comment('班级') VARCHAR(50)" validate:"required"`
	Mobile    string    `form:"mobile" xorm:"not null default '' comment('手机') VARCHAR(11)" validate:"required"`
	CreateAt  time.Time `xorm:"not null created" validate:"lte=128"`
}
