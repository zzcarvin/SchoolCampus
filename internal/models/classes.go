package models

import (
	"time"
)

type Classes struct {
	Id           int       `json:"id" xorm:"not null pk autoincr comment('班级id') INT(5)"`
	DepartmentId int       `json:"department_id" xorm:"comment('院系id') INT(5)"`
	Name         string    `json:"name" xorm:"comment('班级名称') VARCHAR(20)"`
	CreateAt     time.Time `json:"create_at" xorm:"comment('创建时间') DATETIME"`
	ClassType    int       `json:"class_type" xorm:"not null comment('1代表行政班,2为体育班') TINYINT(4)" validate:"required"`
}
