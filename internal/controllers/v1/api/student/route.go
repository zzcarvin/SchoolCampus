package student

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	party.Post("/", create)
	party.Delete("/{id:uint64}", remove)
	party.Put("/{id:uint64}", update)
	party.Get("/{id:uint64}", get)
	party.Get("getcount", getcount)
	party.Get("/get_plan/{id:uint64}", getPlan)
	party.Get("s", search)
	//学生计划进度
	party.Get("/progress", studentPlanProgress)
	//作弊反馈
	party.Get("/studentAbnormal", studentAbnormal)
	//学生每周情况,学生个人页统计数据，比如总共里，跑步有效次数等
	party.Get("/monthProgress", everyMonthProgress)
	party.Get("b", gethistogram)
	//批量导入
	party.Post("/excel", excel)

	party.Post("/update_excel", UpdateExcel)

	//批量修改
	party.Post("/many_update_excel", manyUpdateExcel)

	//各时间段学生完成情况分布图
	//完成次数图表
	party.Get("/completenumber", completenumber)
}
