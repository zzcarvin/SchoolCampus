package student

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party) {

	//获取运动记录
	party.Get("/records", records)
	//获取运动数据
	party.Get("/planprogress/{id:uint64}", planProgress)
	//查询学生
	party.Get("/search", search)

	//获取学生整合信息
	party.Get("/all_info/{id:uint64}", get)

	//学生每周情况,学生个人页统计数据，比如总共里，跑步有效次数等
	party.Get("/week_progress", everyMonthProgress)

}
