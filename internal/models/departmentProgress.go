package models

import "time"

type WeekProgress struct {

	//Id
	Id int `json:"id" xorm:"not null pk autoincr comment('id') INT(10)"`

	//院系id
	DepartmentId int `json:"department_id" xorm:"not null comment('院系id') INT(11)"`

	ClassId int `json:"class_id" xorm:"not null comment('班级id') INT(11)"`

	//每周跑步次数
	Times int `json:"times" xorm:" times not null comment('每周跑步次数') INT(20)"`

	//每周里程
	Distance int `json:"distance" xorm:"distance not null comment('每周里程') INT(30)"`

	//计划百分比
	Percentage float32 `json:"percentage" xorm:"not null comment('计划百分比') FLOAT"`

	Sequence int `json:"sequence" xorm:"not null comment('周序号') tinyint(4)"`

	Tear int `json:"tear" xorm:"term"`

	TermYear int `json:"term_year" xorm:"term_year"`

	Create time.Time `json:"create" xorm:"create_at created"`
}
