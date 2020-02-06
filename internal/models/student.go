package models

import "time"

var StudentInfo Student

type Student struct {
	Id           int       `json:"id" xorm:"autoincr id pk" `
	Year         int       `json:"year" xorm:"year"`
	DepartmentId int       `json:"department_id" xorm:"department_id"`
	UserId       string    `json:"user_id" xorm:"user_id"`
	ClassId      int       `json:"class_id" xorm:"class_id"`
	Name         string    `json:"name" xorm:"name"`
	Gender       int       `json:"gender" xorm:"gender"`
	Code         string    `json:"code" xorm:"code"`
	IdCard       string    `json:"id_card" xorm:"id_card"`
	Cellphone    string    `json:"cellphone" xorm:"cellphone"`
	Continue     int       `json:"continue" xorm:"continue"`
	LastSport    time.Time `json:"last_sport" xorm:"last_sport"`
	CreateAT     time.Time `json:"create_at" xorm:"created 'create_at'"`

	//-- 2019.11.4 新增字段
	PlanId            int `json:"plan_id" xorm:"plan_id"`
	SportDepartmentId int `json:"sport_department_id" xorm:"-"`
	SportClassId      int `json:"sport_class_id" xorm:"-"`
}

type StudentAllInfos struct {
	Id                int       `json:"student_id" xorm:"id"`
	Code              string    `json:"student_code" xorm:"code"`
	Name              string    `json:"name" xorm:"name"`
	Gender            int       `json:"gender" xorm:"gender"`
	ClassName         string    `json:"class_name" xorm:"name"`
	DepartmentName    string    `json:"department_name" xorm:"name"`
	Cellphone         string    `json:"cellphone" xorm:"cellphone"`
	Year              int       `json:"year" xorm:"year"`
	PlanId            int       `json:"plan_id" xorm:"plan_id"`
	CreateAt          time.Time `json:"create_at" xorm:"create_at"`
	DepartmentId      int       `json:"department_id" xorm:"department_id"`
	ClassId           int       `json:"class_id" xorm:"class_id"`
	SportDepartmentId int       `json:"sport_department_id" xorm:"-"`
	SportClassId      int       `json:"sport_class_id" xorm:"-"`
}
type NewStudent struct {
	Department string `json:"department" xorm:"department"`
	Class      string `json:"class" xorm:"class"`   //班级名称
	Name       string `json:"name" xorm:"name"`     //学生名称
	Gender     int    `json:"gender" xorm:"gender"` //性别
	Code       string `json:"code" xorm:"code"`     //学号
}
type StudentId struct {
	ClassId      int `json:"class_id" xorm:"id"` //这是classid
	DepartmentId int `json:"department_id" xorm:"id"`
}
