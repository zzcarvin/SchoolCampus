package student

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/", update)
	party.Get("/{id:uint64}", get)
	party.Get("getcount", getcount)
	party.Get("s", search)
	//学生计划进度
	party.Get("/progress", studentPlanProgress)
	//作弊反馈
	party.Get("/studentAbnormal", studentAbnormal)
	//学生每周情况,学生个人页统计数据，比如总共里，跑步有效次数等
	party.Get("/monthProgress", everyMonthProgress)
	party.Get("b", gethistogram)
	party.Post("execl", execl)

}
