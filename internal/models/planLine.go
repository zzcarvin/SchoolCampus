package models

import "time"

//由两个点组成的路线，防止出现路线中有不能逾越的地理障碍
type PlanLine struct {
	// Required: true
	Id         int         `json:"id" xorm:"autoincr id pk "`
	FenceId    int         `json:"fence_id" xorm:"fence_id"`
	Point1     int         `json:"point1" xorm:"point1"`
	Point1Ptr  *PlanPoints `json:"-" xorm:"-"`
	Point2     int         `json:"point2" xorm:"point2"`
	Point2Ptr  *PlanPoints `json:"-" xorm:"-"`
	CreateAT   time.Time   `json:"create_at" xorm:"created 'create_at'"`
	CoverTimes int         `json:"-" xorm:"-"`
}
