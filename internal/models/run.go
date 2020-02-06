package models

type RequestFinishRun struct {
	PlanStartId  int      `json:"plan_start_id"`
	SchoolId     int      `json:"school_id"`
	PlanId       int      `json:"plan_id"`
	StudentId    int      `json:"student_id"`
	StartTime    string   `json:"start_time"`
	EndTime      string   `json:"end_time"`
	Distance     int      `json:"distance"`
	Duration     int      `json:"duration"`
	Times        int      `json:"times"`
	Calories     float64  `json:"calories"`
	Steps        int      `json:"steps"`
	Pace         int      `json:"pace"`
	Speed        float64  `json:"speed"`
	Paces        []int    `json:"paces"`
	Points       []Points `json:"points"`
	PointsStatus bool     `json:"points_status"`
	//***************************增加男女生配速*************************

	//最大配速
	BoyPaceMax int `json:"boy_pace_max" xorm:"comment('最大配速') INT(11)"`

	//最小配速
	BoyPaceMin int `json:"boy_girl_pace_min" xorm:"comment('最小配速') INT(11)"`
	//最大配速
	GirlPaceMax int `json:"girl_pace_max" xorm:"comment('最大配速') INT(11)"`

	//最小配速
	GirlPaceMin int `json:"girl_pace_min" xorm:"comment('最小配速') INT(11)"`

	//*****************************************************************

	//增加经过点
	PassPoints []Points `json:"pass_points"`
	FaceStatus bool     `json:"face_status" xorm:"-"`

	IBeacon []Points `json:"ibeacon"`
}

type RequestFinishRun1 struct {
	//学生id
	StudentId int `json:"student_id" xorm:"comment('学生id') INT(7)"`
	//计划id
	PlanId int `json:"plan_id" xorm:"comment('计划Id') INT(6)"`
}
