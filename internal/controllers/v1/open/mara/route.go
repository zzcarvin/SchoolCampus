package mara

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//分享
	//party.Get("/records",share)

	//使用app的接口逻辑，放弃旧的逻辑
	party.Post("/mara_apply", ApplyMara)
	party.Get("/mara_status", GetStatus)

}
