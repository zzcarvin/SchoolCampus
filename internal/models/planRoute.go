package models

import "time"

type PlanRoute struct {
	Id        int       `json:"id" xorm:"autoincr"`
	StudentId int       `json:"user_id" xorm:"student_id"`
	PlanId    int       `json:"plan_id" xorm:"plan_id"`
	FenceId   int       `json:"fence_id" xorm:"fence_id"`
	Route     []int64   `json:"route" xorm:"route"`
	CreateAt  time.Time `json:"create_at" xorm:"created"`
}
