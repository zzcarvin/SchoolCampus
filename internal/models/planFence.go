package models

import (
	"time"
)

func (PlanFence) TableName() string {
	return "plan_fence"
}

type PlanFence struct {
	Id        int       `json:"id" xorm:"autoincr 'id' pk"`
	Name      string    `json:"name" xorm:"name"`
	Longitude float64   `json:"longitude" xorm:"longitude"`
	Latitude  float64   `json:"latitude" xorm:"latitude"`
	Points    []Points  `json:"points" xorm:"points"`
	CreateAT  time.Time `json:"create_at" xorm:"created 'create_at'"`
}

type Points struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

//var PlanFence = struct {
//
//	Id       	int   		`json:"id" xorm:"id  pk "`
//	Name   		string   	`json:"name" xorm:"name"`
//	Points 		[]Points 	`json:"points" xorm:"points"`
//	CreateAT    time.Time 	`json:"create_at" xorm:"create_at"`
//	Centre      []Points1	`json:"points1" xorm:"points1"`
//}{}
//
//type Points struct {
//	Lon string `json:"lon" xorm:"lon"`
//	Lat string `json:"lat" xorm:"lat"`
//}
//type Points1 struct {
//	Lon string `json:"lon" xorm:"lon"`
//	Lat string `json:"lat" xorm:"lat"`
//}
//
