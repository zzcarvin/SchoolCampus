package models

import (
	"time"
)

type Department struct {
	Id             int       `json:"id" xorm:"not null pk autoincr comment('院系id') INT(5)"`
	Name           string    `json:"name" xorm:"comment('院系名称') unique VARCHAR(20)"`
	CreateAt       time.Time `json:"create_at" xorm:"DATETIME"`
	DepartmentType int       `json:"department_type" xorm:"not null comment('1为行政系，2为体育系') TINYINT(4)" validate:"required"`
}
