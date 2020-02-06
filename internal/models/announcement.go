package models
import (
	"time"
)

type Announcement struct {
	//Id
	Id int `json:"id" xorm:"not null pk comment('id') INT(11)"`

	//标题
	Title string `json:"title" xorm:"comment('标题') TEXT"`

	//院系
	DepartmentId []int `json:"department_id" xorm:"comment('院系') VARCHAR(50)"`

	//班级
	ClassId []int `json:"class_id" xorm:"comment('班级') VARCHAR(50)"`

	//年级
	GradeId []int `json:"grade_id" xorm:"comment('年级') VARCHAR(50)"`

	//正文
	Content string `json:"content" xorm:"comment('正文') TEXT"`

	//日期
	Date time.Time `json:"date" xorm:"comment('日期') DATE created"`

	//时间
	Time time.Time `json:"time" xorm:"comment('时间') TIME created"`

	//状态
	Status int `json:"status" xorm:"comment('状态') INT(11)"`

	//详情
	Details string `json:"details" xorm:"comment('详情') TEXT"`

	Apptime string `json:"apptime" xorm:"comment('时间') VARCHAR(100)"`


}

