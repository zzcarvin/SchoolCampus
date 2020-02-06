package models

import "time"

type User struct {
	// Required: true
	Id         int       `json:"id"     xorm:"id"    validate:"required,lte=7"`
	Code       string    `json:"code"   xorm:"code"  validate:"lte=128"`
	Nickname   string    `json:"nickname" xorm:"nickname" validate:"lte=128"`
	Password   string    `json:"password" xorm:"password" validate:"lte=128"`
	Cellphone  string    `json:"cellphone" xorm:"cellphone" validate:"lte=11"`
	IdCard     string    `json:"id_card" xorm:"id_card" validate:"lte=128"`
	Gender     int       `json:"gender" xorm:"gender default(1)"  validate:"lte=128"`
	Face       string    `json:"face" xorm:"face"  validate:"lte=128"`
	FaceStatus int       `json:"face_status" xorm:"face_status" validate:"let=1"`
	Avatar     string    `json:"avatar" xorm:"avatar" validate:"lte=128"`
	CreateAt   time.Time `json:"create_at" xorm:"created" validate:"lte=128"`
}
