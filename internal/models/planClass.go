package models

type PlanClass struct {
	Id           int `json:"id" xorm:"id"`
	PlanId       int `json:"plan_id" xorm:"plan_id"`
	ClassId      int `json:"class_id" xorm:"class_id"`
	DepartmentId int `json:"department_id" xorm:"department_id"`
	Gender       int `json:"gender" xorm:"gender"`
}
