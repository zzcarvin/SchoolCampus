package models

type StudentClass struct {
	Id           int `json:"id" xorm:"id"`
	StudentId    int `json:"student_id" xorm:"student_id"`
	DepartmentId int `json:"department_id" xorm:"department_id"`
	ClassId      int `json:"class_id" xorm:"class_id"`
}
