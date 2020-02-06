package models

import "time"

type Plan struct {
	//计划id
	Id int `json:"id" xorm:"not null pk autoincr comment('计划Id') INT(11)"`

	//院系id
	DepartmentId int `json:"department_id" xorm:"not null comment('院系id') index INT(5)"`

	//班级id
	ClassId int `json:"class_id" xorm:"not null comment('班级id') index INT(5)"`

	//计划名称
	Name string `json:"name" xorm:"not null comment('计划名称') VARCHAR(255)"`

	//计划开始日期
	DateBegin time.Time `json:"date_begin" xorm:"not null comment('计划开始日期') DATETIME"`

	//计划结束日期
	DateEnd time.Time `json:"date_end" xorm:"not null comment('计划结束日期') DATETIME"`

	//总里程
	TotalDistance int `json:"total_distance" xorm:"not null comment('总里程') INT(4)"`

	//总次数
	TotalTimes int `json:"total_times" xorm:"not null comment('总次数') INT(4)"`

	//最低周跑次数
	MinWeekTimes int `json:"min_week_times" xorm:"not null comment('最低周跑次数') INT(2)"`

	//最低周跑里程
	MinWeekDistance int `json:"min_week_distance" xorm:"not null comment('最低周跑里程') INT(11)"`

	//最低单次里程
	MinSingleDistance int `json:"min_single_distance" xorm:"not null comment('最低单次里程') INT(11)"`

	//最高周跑次数
	MaxWeekTimes int `json:"max_week_times" xorm:"not null comment('最高周跑次数') INT(11)"`

	//最高周跑里程
	MaxWeekDistance int `json:"max_week_distance" xorm:"not null comment('最高周跑里程') INT(11)"`

	//最高单次里程
	MaxSingleDistance int `json:"max_single_distance" xorm:"not null comment('最高单次里程') INT(11)"`

	//最高日跑次数
	MaxDayTimes int `json:"max_day_times" xorm:"not null comment('最高日跑次数') INT(2)"`

	//必跑日
	MustRunDay string `json:"must_run_day" xorm:"not null comment('必跑日') VARCHAR(20)"`

	//最小配速(单位s)
	MinPace int `json:"min_pace" xorm:"not null comment('最小配速(单位s)') INT(11)"`

	//最大配速(单位s)
	MaxPace int `json:"max_pace" xorm:"not null comment('最大配速(单位s)') INT(11)"`

	//最短时间(单位s)
	MinTimeLong int `json:"min_time_long" xorm:"not null comment('最短时间(单位s)') INT(11)"`

	//最长时间(单位s)
	MaxTimeLong int `json:"max_time_long" xorm:"not null comment('最长时间(单位s)') INT(11)"`

	//创建时间
	CreateTime time.Time `json:"create_time" xorm:"not null comment('创建时间') DATETIME"`

	//更新时间
	UpdateTime time.Time `json:"update_time" xorm:"not null comment('更新时间') DATETIME"`

	//计划终止状态，1：计划有效中，2：计划已终止
	Stop int `json:"stop" xorm:"not null comment('计划终止状态，1：计划有效中，2：计划已终止') index TINYINT(1)"`

	//年级
	Year int `json:"year" xorm:"not null comment('年级') INT(11)"`

	//学期 1 第一学期 2 第二学期
	Term int `json:"term" xorm:"not null comment('学期 1 第一学期 2 第二学期') TINYINT(1)"`

	//学年
	TermYear int `json:"term_year" xorm:"not null comment('学年') INT(11)"`

	//1为行政班，2为运动班
	PlanType int `json:"plan_type" xorm:"not null comment('1为行政班，2为运动班') INT(11)"`
	//****************增加性别字段
	Gender    int             `json:"gender" xorm:"gender"`
	TimeFrame []PlanTimeFrame `json:"time_frame" xorm:"-"`
}

type PlanTimeFrame struct {
	Id            int       `json:"id" xorm:"not null pk autoincr comment('计划Id') INT(11)"`
	PlanId        int       `xorm:"not null INT(11)"`
	DurationBegin string    `xorm:"not null comment('每日计划开始时间') VARCHAR(5)"`
	DurationEnd   string    `xorm:"not null comment('每日计划结束时间') VARCHAR(5)"`
	CreateAt      time.Time `json:"create_at" xorm:"create_at created"`
}

type PlanAndProgress struct {
	PlanId                   int    `json:"plan_id" xorm:"id"`
	Name                     string `json:"name" xorm:"name"`             //计划名称
	Date                     string `json:"date" xorm:"date"`             //计划日期
	DateBegin                string `json:"date_begin" xorm:"date_begin"` //计划开始日期
	DateEnd                  string `json:"date_end" xorm:"date_end"`
	Types                    int    `json:"types" xorm:"types"`
	BoyTotalDistance         int    `json:"boy_total_distance" xorm:"boy_total_distance"` //计划总里程
	BoySingleMindistance     int    `json:"boy_single_mindistance" xorm:"boy_single_mindistance"`
	GirlTotalDistance        int    `json:"girl_total_distance" xorm:"girl_total_distance"`
	GirlSingleMindistance    int    `json:"girl_single_mindistance" xorm:"girl_single_mindistance"`
	BoyWeekDistance          int    `json:"boy_week_distance" xorm:"boy_week_distance"`
	BoyWeekSingleMaxdistance int    `json:"boy_week_single_maxdistance" xorm:"boy_week_single_maxdistance"`
	GirlWeekDistance         int    `json:"girl_week_distance" xorm:"girl_week_distance"`

	GirlWeekSingleMaxdistance int     `json:"girl_week_single_maxdistance" xorm:"girl_week_single_maxdistance"`
	BoyTotalTimes             int     `json:"boy_total_times" xorm:"boy_total_times"`
	BoyTimesMindistance       int     `json:"boy_times_mindistance" xorm:"boy_times_mindistance"`
	GirlTotalTimes            int     `json:"girl_total_times" xorm:"girl_times_times"`
	GirlTimesMindistance      int     `json:"girl_times_mindistance" xorm:"girl_times_mindistance"`
	Gender                    int     `json:"gender" xorm:"gender"`
	Duration                  string  `json:"duration" xorm:"duration"` //计划时间段
	DurationBegin             string  `json:"duration_begin" xorm:"duration_begin"`
	DurationEnd               string  `json:"duration_end" xorm:"duration_end"`
	StrideFrequency           float32 `json:"stride_frequency" xorm:"stride_frequency"` //计划步频
	Pace                      int     `json:"pace" xorm:"pace"`                         //计划配速
	PlanProgressId            int     `json:"plan_progress_id" xorm:"id"`
	Distance                  int     `json:"distance" xorm:"distance"`         //进度总公里数
	ProgressDuration          int     `json:"progressDuration" xorm:"duration"` //进度总时间
	Times                     int     `json:"times" xorm:"times"`               //进度总次数
	Calories                  int     `json:"calories" xorm:"calories"`         //进度总卡路里
	Steps                     int     `json:"steps" xorm:"steps"`               //进度总步数
	Weeks                     int     `json:"weeks"`
}

//type PlanAndProgress struct {
//	PlanProgress
//	Plan
//}
