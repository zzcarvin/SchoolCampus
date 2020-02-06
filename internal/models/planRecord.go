package models

import "time"

type PlanRecord struct {

	//计划完成记录id
	Id int `json:"id" xorm:"not null pk autoincr comment('计划完成记录id') INT(7)"`
	//学校id
	SchoolId int `json:"school_id" xorm:"comment('学校id') INT(5)"`

	//计划id
	PlanId int `json:"plan_id" xorm:"comment('计划Id') INT(6)"`

	//学生id
	StudentId int `json:"student_id" xorm:"comment('学生id') INT(7)"`

	//类型，0正常记录，1跑步开始时间超过24小时记录
	Type int `json:"type" xorm:"comment('类型，0正常记录，1跑步开始时间超过24小时记录') INT(2)"`

	//开始时间
	StartTime time.Time `json:"start_time" xorm:"comment('开始时间') DATETIME"`

	//结束时间
	EndTime time.Time `json:"end_time" xorm:"comment('结束时间') DATETIME"`

	//里程（米）
	Distance int `json:"distance" xorm:"comment('里程（米）') INT(3)"`

	//运动时间
	Duration int `json:"duration" xorm:"comment('运动时间') INT(10)"`

	//卡路里（卡）
	Calories float64 `json:"calories" xorm:"comment('卡路里（卡）') FLOAT"`

	//配速
	Pace int `json:"pace" xorm:"comment('配速') VARCHAR(255)"`

	//格式化配速
	FormPace string `json:"form_pace" xorm:"form_pace"`

	//步数
	Steps int `json:"steps" xorm:"comment('步数') INT(5)"`

	//平均步频
	Frequency float64 `json:"frequency" xorm:"frequency"`

	//步频数组
	Frequencies []int `json:"frequencies" xorm:"frequencies comment('步频数组') TEXT"`

	//点
	//对于10118byte，数组内元素（也就是一个完整的坐标）大概是206个，text格式最多存65535byte，以10118*6算为60780byte，安全起见，允许存入最多元素为1200个
	Points []Points `json:"points" xorm:"comment('点') TEXT" validate:"required|maxLen:1200"`

	//创建时间
	CreateAt time.Time `json:"create_at" xorm:"comment('创建时间') DATETIME created"`

	//
	Times float64 `json:"times" xorm:"INT(11)"`

	//速度
	Speed float64 `json:"speed" xorm:"comment('速度') DOUBLE"`

	//状态
	Status int `json:"status" xorm:"comment('状态') INT(11)"`

	//无效码
	InvalidCode []int `json:"invalid_code" xorm:"invalid_code"`
	//性别
	Gender int `json:"gender" xorm:"gender"`
	//院系
	DepartmentId int `json:"department_id" xorm:"department_id"`

	//系统分配的路径id
	RouteId int `json:"route_id" xorm:"route_id"`

	//经过的蓝牙点id
	PassPoints []int `json:"pass_points" xorm:"pass_points"`

	PlanRecordString []string `json:"planrecordstring"  xorm:"-"`
}
