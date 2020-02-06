package models

import (
	"time"
)

type Account struct {

	// required: true

	Id        int       `json:"id"            xorm:"id pk autoincr"   validate:"required,lte=128"`
	Username  string    `json:"username"      xorm:"username"       validate:"lte=128"`
	Password  string    `json:"password"      xorm:"password"       validate:"lte=128"`
	Nickname  string    `json:"nickname"      xorm:"nickname"       validate:"lte=128"`
	RoleId    int       `json:"role_id"       xorm:"role_id"        validate:"lte=128"`
	CreateAt  time.Time `json:"create_at"     xorm:"created"      validate:"lte=128"`
	Lastlogin time.Time `json:"last_login" xorm:"last_login updated"  validate:"lte=128"`
}
type Password struct {
	//Id          int    `json:"id"            xorm:"id pk"             validate:"required,lte=128"`
	Oldpassword string `json:"oldpassword"      xorm:"—"       validate:"lte=128"`
	Newpassword string `json:"newpassword"      xorm:"password"       validate:"lte=128"`
}
