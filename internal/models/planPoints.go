package models

import "time"

type PlanPoints struct {
	// Required: true
	Id        int         `json:"id" xorm:"autoincr pk id"`
	Name      string      `json:"name" xorm:"name"`
	FenceId   int         `json:"fence_id" xorm:"fence_id"`
	Longitude float64     `json:"longitude" xorm:"longitude DOUBLE"`
	Latitude  float64     `json:"latitude" xorm:"latitude DOUBLE"`
	Type      int         `json:"type" xorm:"type"`
	Uuid      string      `json:"uuid" xorm:"uuid"`
	Lines     []*PlanLine `json:"-" xorm:"-"` //插入点时忽略该字段
	CreateAT  time.Time   `json:"create_at" xorm:"created 'create_at'"`
	Deleted   int         `json:"-" xorm:"'deleted'"`
	DeleteAt  time.Time   `json:"delete_at" xorm:"created delete_at"`
}

type IBeaconsPoints32 struct {
	Id        int       `json:"id" xorm:"autoincr pk id"`
	Name      string    `json:"name" xorm:"name"`
	FenceId   int       `json:"fence_id" xorm:"fence_id"`
	Longitude float32   `json:"longitude" xorm:"longitude DOUBLE"`
	Latitude  float32   `json:"latitude" xorm:"latitude DOUBLE"`
	Type      int       `json:"type" xorm:"type"`
	Uuid      string    `json:"uuid" xorm:"uuid"`
	CreateAT  time.Time `json:"create_at" xorm:"created 'create_at'"`
}
