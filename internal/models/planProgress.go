package models

import "time"

type PlanProgress struct {
	// Required: true
	Id           int       `json:"id" xorm:"autoincr id pk" `
	DepartmentId int       `json:"department_id" xorm:"department_id"`
	StudentId    int       `json:"student_id" xorm:"student_id"`
	PlanId       int       `json:"plan_id" xorm:"plan_id"`
	Distance     int       `json:"distance" xorm:"distance"`
	Duration     int       `json:"duration" xorm:"duration"`
	Times        int       `json:"times" xorm:"times"`
	Calories     float64   `json:"calories" xorm:"calories"`
	Steps        int       `json:"steps" xorm:"steps"`
	Status       int       `json:"status" xorm:"status"`
	CreateAt     time.Time `json:"create_at" xorm:"create_at created"`
	Year         int       `json:"year" xorm:"year"`

	//每周有效跑量
	WeekDistance int `json:"week_distance" xorm:"not null comment('每周有效跑量') INT(11)"`

	//每周有效跑步次数
	WeekTimes int `json:"week_times" xorm:"not null comment('每周有效跑步次数') INT(11)"`

	//完成进度
	CompleteProgress float32 `json:"complete_progress" xorm:"comment('完成进度') FLOAT(32)"`

	//周完成进度
	WeekCompleteProgress string `json:"week_complete_progress" xorm:"comment('周完成进度') VARCHAR(200)"`

	//班级
	ClassId int `json:"class_id" xorm:"not null comment('班级') INT(11)"`
}
