package record

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Get("/{id:uint64}", get) //获取单条记录
	//party.Get("/search", search)
	//新的查询单个学生所有运动
	party.Get("/search", newSearch)
	party.Get("/best/{id:uint64}", best)
	party.Post("/duration/{id:uint64}", duration)
	party.Post("/newduration", newDuration)

	//分享
	party.Get("/share", share)

	//打卡点
	party.Get("/points", passPoints)

}
