package models

import "time"

type Teacher struct {
	Id           int       `json:"id" xorm:"autoincr id pk" `
	SchoolId     int       `json:"school_id" xorm:"school_id"`
	DepartmentId int       `json:"department_id" xorm:"department_id"`
	UserId       string    `json:"user_id" xorm:"user_id"`
	Name         string    `json:"name" xorm:"name"`
	Gender       int       `json:"gender" xorm:"gender"`
	Code         string    `json:"code" xorm:"code"`
	IdCard       string    `json:"id_card" xorm:"id_card"`
	Cellphone    string    `json:"cellphone" xorm:"cellphone"`
	CreateAT     time.Time `json:"create_at" xorm:"created 'create_at'"`
}

type TeacherAllInfos struct {
	Id             int       `json:"teacher_id" xorm:"id"`
	Code           string    `json:"code" xorm:"code"`
	Name           string    `json:"name" xorm:"name"`
	Gender         int       `json:"gender" xorm:"gender"`
	DepartmentName string    `json:"department_name" xorm:"name"`
	Cellphone      string    `json:"cellphone" xorm:"cellphone"`
	CreateAt       time.Time `json:"create_at" xorm:"create_at"`
	DepartmentId   int       `json:"department_id" xorm:"department_id"`
}
