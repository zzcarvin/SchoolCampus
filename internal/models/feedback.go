package models

import "time"

type Feedback struct {
	Id        int    `json:"id" xorm:"id pk autoincr"`
	StudentId int    `json:"student_id" xorm:"student_id"`
	Content   string `json:"content" xorm:"content"`
	RecordId  int    `json:"record_id" xorm:"record_id"`
	Status    int    `json:"status" xorm:"status"`
	//详情
	ProblemTypes  []string  `json:"problem_types" xorm:"comment('问题反馈类型') TEXT"`
	CreateAt      time.Time `json:"create_at"     xorm:"created"      validate:"lte=128"`
	Reply_message string    `json:"reply_message" xorm:"reply_message"`
	Check_status  int       `json:"feedback_status" xorm:"Check_status"`
}
